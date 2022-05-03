
#include <ESP8266WiFi.h>
#include <MFRC522.h>
#include <PubSubClient.h>
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


bool beep_on = false;
long beep_activated = 0;

MFRC522 mfrc522(SS_PIN, RST_PIN);
MFRC522::MIFARE_Key key;

WiFiClient espClient;
PubSubClient client(espClient);

String ssid = "";
String psk = "";

Config config;





void safeWrite(byte* data, int size);
void reconnect_wifi();





void publishFlashStringHelper(const char* topic, const __FlashStringHelper *text) {
  int length = 0;
  length += strlen_P((const char*)text);
  char buffer[length+1];

  strncpy_P(buffer, (const char*)text, length);  
  buffer[length] = '\0';

  client.publish(topic, buffer);
}





void copy_byte_to_cstr(byte *from, int from_length, char* to, int to_length) {
  if (from_length >= to_length) {
    from_length = to_length - 1;
  }

  int i = 0;
  for (i = 0; i < from_length; i++) {
    to[i] = from[i];
  }
  to[i] = '\0';
}



void onDoorlockWrite(char *topic, byte *payload, unsigned int length) {
  int decoded_length = 128;
  unsigned char buffer[decoded_length];

  char c_str_payload[length+1];
  copy_byte_to_cstr(payload, length, c_str_payload, length+1);
  
  decoded_length = decode_base64((unsigned char*)c_str_payload, buffer);

  safeWrite(buffer, decoded_length);
}

void onConnectionEstablished() {
  client.publish("doorlock/" CHIP_ID "/status", "true", true);
  client.subscribe("doorlock/" CHIP_ID "/write/data");
}


void callback(char *topic, byte* payload, unsigned int length) {
  if (strcmp(topic, "doorlock/" CHIP_ID "/write/data") == 0) {
    onDoorlockWrite(topic, payload, length);
  }
}





void setup_wifi() {

  strcmp(config.server, "<SomeIP>");
  config.port = 1883;
  strcmp(config.user, "User");
  strcmp(config.password, "Password");

  WiFi.setSleepMode(WIFI_NONE_SLEEP);
  WiFi.persistent(false);

  reconnect_wifi();
}

void reconnect_wifi() {

  if (psk == "" || ssid == "") {
    setup_wifi();
    return;
  }

  if (WiFi.status() == WL_CONNECTED) {
    return;
  }
  
  
  // Connect to wifi client
  delay(10);
  // We start by connecting to a WiFi network
  Serial.println();
  Serial.print("Connecting to ");
  Serial.println(ssid);

  WiFi.mode(WIFI_STA);
  WiFi.begin(ssid, psk);

  while (WiFi.status() != WL_CONNECTED) {
    delay(500);
    Serial.print(".");
  }

  randomSeed(micros());

  Serial.println("");
  Serial.println("WiFi connected");
  Serial.println("IP address: ");
  Serial.println(WiFi.localIP());
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


  SPI.begin();
  mfrc522.PCD_Init();

  initRfidKey();

  setup_wifi();

  client.setServer("<SomeIP>", 1883);
  client.setCallback(callback);
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
    client.publish("doorlock/" CHIP_ID "/error", "To many bytes. Max size is 48 bytes.");
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
    client.publish("doorlock/" CHIP_ID "/error", "Card not present anymore");
    return;
  }

  MFRC522::StatusCode status;
  mfrc522.PCD_StopCrypto1();
  status = (MFRC522::StatusCode) mfrc522.PCD_Authenticate(MFRC522::PICC_CMD_MF_AUTH_KEY_A, TRAILING_SECTOR, &key, &(mfrc522.uid));
  if (status != MFRC522::STATUS_OK) {
    publishFlashStringHelper("doorlock/" CHIP_ID "/error", mfrc522.GetStatusCodeName(status));
    return;
  }

  int blockAddr = DEFAULT_BLOCK;


  // Write 3 Blocks, 16 Bytes each. buffer will always be 48 Bytes in length.
  bool success = true;
  status = (MFRC522::StatusCode) mfrc522.MIFARE_Write(blockAddr, buffer, 16);
  if (status != MFRC522::STATUS_OK) {
    publishFlashStringHelper("doorlock/" CHIP_ID "/error", mfrc522.GetStatusCodeName(status));
    success = false;
  }


  status = (MFRC522::StatusCode) mfrc522.MIFARE_Write(blockAddr + 1, buffer + 16, 16);
  if (status != MFRC522::STATUS_OK) {
    publishFlashStringHelper("doorlock/" CHIP_ID "/error", mfrc522.GetStatusCodeName(status));
    success = false;
  }


  status = (MFRC522::StatusCode) mfrc522.MIFARE_Write(blockAddr + 2, buffer + 32, 16);
  if (status != MFRC522::STATUS_OK) {
    publishFlashStringHelper("doorlock/" CHIP_ID "/error", mfrc522.GetStatusCodeName(status));
    success = false;
  }

  // If success, everything worked.
  if (!success) {
    client.publish("doorlock/" CHIP_ID "/error", "Data could not be written onto chip.");
    client.publish("doorlock/" CHIP_ID "/write/ok", "false");
  } else {
    client.publish("doorlock/" CHIP_ID "/write/ok", "true");
  }

  // Halt PICC
  mfrc522.PICC_HaltA();
  // Stop encryption on PCD
  mfrc522.PCD_StopCrypto1();
}



void rfidRead(byte* data, int size) {
  if (size > MAX_BYTES) {
    client.publish("doorlock/" CHIP_ID "/error", "To many bytes. Max size is 48 bytes.");
    return;
  }

  MFRC522::StatusCode status;


  status = (MFRC522::StatusCode) mfrc522.PCD_Authenticate(MFRC522::PICC_CMD_MF_AUTH_KEY_A, TRAILING_SECTOR, &key, &(mfrc522.uid));
  if (status != MFRC522::STATUS_OK) {
      publishFlashStringHelper("doorlock/" CHIP_ID "/error", mfrc522.GetStatusCodeName(status));
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
      publishFlashStringHelper("doorlock/" CHIP_ID "/error", mfrc522.GetStatusCodeName(status));
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
  beep_on = false;
}

void beepOn() {
  digitalWrite(BEEP_PIN, HIGH);
  beep_on = true;
  beep_activated = millis();
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
  client.publish("doorlock/" CHIP_ID "/read", bufferString.c_str());

  // Halt PICC
  mfrc522.PICC_HaltA();
  // Stop encryption on PCD
  mfrc522.PCD_StopCrypto1();

  beepOn();
}



void reconnect_mqtt() {

  while(!client.connected()) {
    Serial.print("Attempting MQTT connection...");
    // Create a random client ID
    String clientId = "ESP8266Client-";
    clientId += String(random(0xffff), HEX);
    // Attempt to connect
    if (client.connect(clientId.c_str(), "doorlock/" CHIP_ID "/opener/status", 0, true, "false")) {
      Serial.println("connected");
      onConnectionEstablished();
    } else {
      reconnect_wifi();
      Serial.print("failed, rc=");
      Serial.print(client.state());
      Serial.println(" try again in 5 seconds");
      delay(5000);
    }
  }
}




void loop() {

  reconnect_wifi();

  if(beep_on && millis() - beep_activated >= 200) {
    beepOff();
  }

  if(!client.connected()) {
    reconnect_mqtt();
  }

  client.loop();
  rfidLoop();
  delay(100);

}
