
#include <WiFiManager.h>
#include <EEPROM.h>
#include <MFRC522.h>
#include <EspMQTTClient.h>
#include <base64.hpp>

#define RST_PIN 2
#define SS_PIN 15
#define BEEP_PIN 5
#define PORTAL_PIN 0

#define MAX_BYTES 48
#define DEFAULT_BLOCK 4 // This will start at Sector 0, Block 0
#define TRAILING_SECTOR 7 // Authentication sector.

#define CHIP_ID "02"

#define MAX_TIMEOUT 20000


typedef struct {
  char server[40] = "";
  char user[40] = "";
  char password[40] = "";
  int port = 0;
} Config;


MFRC522 mfrc522(SS_PIN, RST_PIN);
MFRC522::MIFARE_Key key;

EspMQTTClient *client = NULL;
String ssid = "";
String psk = "";
String willTopic;
Config config;

WiFiManager manager;

String last_uid = "";


bool shouldSave = false;

void onSaveCallback() {
  shouldSave = true;
}

// Init default key. This may change during development
void initRfidKey() {
  for (byte i = 0; i < 6; i++) {
    key.keyByte[i] = 0xFF;
  }
}

void setup() {

  pinMode(PORTAL_PIN, INPUT);

  pinMode(BEEP_PIN, OUTPUT);
  digitalWrite(BEEP_PIN, HIGH);
  delay(200);
  digitalWrite(BEEP_PIN, LOW);

  Serial.begin(9600);



  connect(true);


  SPI.begin();
  mfrc522.PCD_Init();

  initRfidKey();
}


void connect(bool auto_connect) {
  // Try connect to wifi and or mqtt.
  char port[6] = "";
  WiFiManagerParameter mqtt_server("server", "mqtt server", config.server, 40);
  WiFiManagerParameter mqtt_password("password", "mqtt password", config.password, 40);
  WiFiManagerParameter mqtt_user("user", "mqtt user", config.password, 40);
  WiFiManagerParameter mqtt_port("port", "mqtt port", port, 6);
  manager.setSaveConfigCallback(onSaveCallback);
  manager.setConnectTimeout(60);
  manager.addParameter(&mqtt_server);
  manager.addParameter(&mqtt_port);
  manager.addParameter(&mqtt_user);
  manager.addParameter(&mqtt_password);

  
    int result = 0;
    if (auto_connect) {
      result = manager.autoConnect("ConfigAP");
    } else {
      result = manager.startConfigPortal();
    }

  config.port = atoi(mqtt_port.getValue());
  strcpy(config.server, mqtt_server.getValue());
  strcpy(config.user, mqtt_user.getValue());
  strcpy(config.password, mqtt_password.getValue());


  // Save mqtt data
  EEPROM.begin(sizeof(config));
  if (shouldSave && result) {
    Serial.println("Should save");
    EEPROM.put(0, config);
    EEPROM.commit();
  } else {
    EEPROM.get(0, config);
  }
  EEPROM.end();


  psk = WiFi.psk();
  ssid = WiFi.SSID();
  WiFi.disconnect();
  WiFi.setSleepMode(WIFI_NONE_SLEEP);
  WiFi.setAutoReconnect(true);
  WiFi.persistent(true);



  client = new EspMQTTClient(
    ssid.c_str(),
    psk.c_str(),
    config.server,
    config.user,
    config.password,
    "chipread:" CHIP_ID,
    config.port
  );
  client->setKeepAlive(15);

  client->enableDebuggingMessages();
  willTopic = "doorlock/";
  willTopic += CHIP_ID;
  willTopic += "/status";
  client->enableLastWillMessage("doorlock/" CHIP_ID "/status", "false", true);
  client->setMqttReconnectionAttemptDelay(5000);
  client->setWifiReconnectionAttemptDelay(60000);
}

bool reselect_card() {
  //-------------------------------------------------------
  // Can also be used to see if card still available,
  // true means it is false means card isnt there anymore
  //-------------------------------------------------------
  byte s;
  byte req_buff[2];
  byte req_buff_size = 2;
  mfrc522.PCD_StopCrypto1();
  s = mfrc522.PICC_HaltA();
  delay(50);
  s = mfrc522.PICC_WakeupA(req_buff, &req_buff_size);
  dump_byte_array(req_buff, req_buff_size);
  delay(50);
  s = mfrc522.PICC_Select( &(mfrc522.uid), 0);
  if ( mfrc522.GetStatusCodeName((MFRC522::StatusCode)s) == F("Timeout in communication.") ) {
    return false;
  }
  return true;
}


bool compareByteArray(byte* a, int sizeA, byte* b, int sizeB) {
  if (sizeA != sizeB) {
    return false;
  }
  for (int i = 0; i < sizeA; i++) {
    if (a[i] != b[i]) {
      return false;
    }
  }
  return true;
}


void safeWrite(byte* data, int size) {
  if (size > MAX_BYTES) {
    client->publish("doorlock/" CHIP_ID "/error", "To many bytes. Max size is 48 bytes.");
    last_uid = "";
    return;
  }

  byte buffer[MAX_BYTES];
  for (int i = 0; i < MAX_BYTES; i++) {
    buffer[i] = 0;
  }
  for (int i = 0; i < size; i++) {
    buffer[i] = data[i];
  }

  // Hardreset mfrc522 to write data. At this point, it would be a coind flip if mfrc522 is ready
  // to write data on already placed chip. In order to remove randomness, hard reset.
  if (!reselect_card()) {
    client->publish("doorlock/" CHIP_ID "/error", "Card not present anymore");
    last_uid = "";
    return;
  }

  MFRC522::StatusCode status;
  mfrc522.PCD_StopCrypto1();
  status = (MFRC522::StatusCode) mfrc522.PCD_Authenticate(MFRC522::PICC_CMD_MF_AUTH_KEY_A, TRAILING_SECTOR, &key, &(mfrc522.uid));
  if (status != MFRC522::STATUS_OK) {
    client->publish("doorlock/" CHIP_ID "/error", mfrc522.GetStatusCodeName(status));
    last_uid = "";
    return;
  }

  int blockAddr = DEFAULT_BLOCK;


  // Write 3 Blocks, 16 Bytes each. buffer will always be 48 Bytes in length.
  bool success = true;
  status = (MFRC522::StatusCode) mfrc522.MIFARE_Write(blockAddr, buffer, 16);
  if (status != MFRC522::STATUS_OK) {
    client->publish("doorlock/" CHIP_ID "/error", mfrc522.GetStatusCodeName(status));
    success = false;
  }


  status = (MFRC522::StatusCode) mfrc522.MIFARE_Write(blockAddr + 1, buffer + 16, 16);
  if (status != MFRC522::STATUS_OK) {
    client->publish("doorlock/" CHIP_ID "/error", mfrc522.GetStatusCodeName(status));
    success = false;
  }


  status = (MFRC522::StatusCode) mfrc522.MIFARE_Write(blockAddr + 2, buffer + 32, 16);
  if (status != MFRC522::STATUS_OK) {
    client->publish("doorlock/" CHIP_ID "/error", mfrc522.GetStatusCodeName(status));
    success = false;
  }

  // If success, everything worked.
  if (!success) {
    client->publish("doorlock/" CHIP_ID "/error", "Data could not be written onto chip.");
    client->publish("doorlock/" CHIP_ID "/write/ok", "false");
  } else {
    client->publish("doorlock/" CHIP_ID "/write/ok", "true");
  }

  // Halt PICC
  mfrc522.PICC_HaltA();
  // Stop encryption on PCD
  mfrc522.PCD_StopCrypto1();
  last_uid = "";
}

void onDoorlockWrite(const String &msg) {
  int decoded_length = 128;
  unsigned char buffer[decoded_length];
  decoded_length = decode_base64((unsigned char*)msg.c_str(), buffer);

  safeWrite(buffer, decoded_length);
}


void onConnectionEstablished() {
  client->publish("doorlock/" CHIP_ID "/status", "true", true);
  client->subscribe("doorlock/" CHIP_ID "/write/data", onDoorlockWrite);
}


void rfidRead(byte* data, int size) {
  if (size > MAX_BYTES) {
    client->publish("doorlock/" CHIP_ID "/error", "To many bytes. Max size is 48 bytes.");
    return;
  }

  MFRC522::StatusCode status;


  status = (MFRC522::StatusCode) mfrc522.PCD_Authenticate(MFRC522::PICC_CMD_MF_AUTH_KEY_A, TRAILING_SECTOR, &key, &(mfrc522.uid));
  if (status != MFRC522::STATUS_OK) {
    client->publish("doorlock/" CHIP_ID "/error", mfrc522.GetStatusCodeName(status));
    return;
  }

  byte buffer[18] = {0};
  byte bufferSize = 18;
  int offset = 0;
  int blockAddr = DEFAULT_BLOCK;
  for (int i = 0; i < 3; i++, blockAddr++) {
    if (blockAddr + 1 % 4 == 0) {
      blockAddr++;
    }
    int remaining = min(16, size - offset);
    if (remaining == 0) {
      break;
    }

    status = (MFRC522::StatusCode) mfrc522.MIFARE_Read(blockAddr, buffer, &bufferSize);
    if (status != MFRC522::STATUS_OK) {
      client->publish("doorlock/" CHIP_ID "/error", mfrc522.GetStatusCodeName(status));
      return;
    }

    for (int j = 0; j < remaining; j++) {
      data[j + offset] = buffer[j];
    }

    offset += remaining;

  }

}

void beepOff() {
  digitalWrite(BEEP_PIN, LOW);
}

void rfidLoop() {

  if ( ! mfrc522.PICC_IsNewCardPresent()) {
    return;
  }

  // Select one of the cards
  if ( ! mfrc522.PICC_ReadCardSerial())
    return;
  size_t base64_size = 127;
  byte base64_uid[base64_size + 1];
  size_t encoded = encode_base64((unsigned char*)mfrc522.uid.uidByte, mfrc522.uid.size, (unsigned char*)base64_uid);

  byte buffer[MAX_BYTES];
  rfidRead(buffer, MAX_BYTES); 
  base64_size = 127;
  byte base64_data[base64_size + 1];
  encoded = encode_base64((unsigned char*)buffer, MAX_BYTES, (unsigned char*) base64_data);

  String bufferString = "{ \"uid\": \"";
  bufferString += String((char *) base64_uid);
  bufferString += String("\", \"data\": \"");
  bufferString += String((char *)base64_data);
  bufferString += "\" }";
  client->publish("doorlock/" CHIP_ID "/read", bufferString);

  last_uid = String((char*) base64_uid);
  // Halt PICC
  mfrc522.PICC_HaltA();
  // Stop encryption on PCD
  mfrc522.PCD_StopCrypto1();

  digitalWrite(BEEP_PIN, HIGH);
  client->executeDelayed(200, beepOff);
}

void loop() {

  if(digitalRead(PORTAL_PIN) == LOW) {
    connect(false);
  }

  client->loop();
  rfidLoop();
  delay(200);

}

// Helper routine to dump a byte array as hex values to Serial
void dump_byte_array(byte *buffer, byte bufferSize) {
  Serial.printf("Size: %d\n", bufferSize);
  for (byte i = 0; i < bufferSize; i++) {
    Serial.printf("%02X ", buffer[i]);
  }
  Serial.printf("\n");
}
