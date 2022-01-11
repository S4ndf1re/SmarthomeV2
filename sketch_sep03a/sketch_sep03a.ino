/*
serialRelayTest.ino
*/

#include <WiFiManager.h>
#include <EEPROM.h>
#include <EspMQTTClient.h>


byte relON[] = {0xA0, 0x01, 0x01, 0xA2};
byte relOFF[] = {0xA0, 0x01, 0x00, 0xA1};


#define CHIP_ID "01"
#define MAX_TIMEOUT 20000

typedef struct {
  char server[40] = "";
  char user[40] = "";
  char password[40] = "";
  int port = 0;
} Config;

EspMQTTClient *client = NULL;
String ssid = "";
String psk = "";
String willTopic;
Config config;

WiFiManager manager;

bool shouldSave = false;

void onSaveCallback() {
  shouldSave = true;
}


void setup(void) {
  Serial.begin(9600);
  
  Serial.write(relOFF, sizeof(relOFF));
  delay(10);
  Serial.write(relOFF, sizeof(relOFF));

  // Try connect to wifi and or mqtt.
  char port[6] = "";
  WiFiManagerParameter mqtt_server("server", "mqtt server", config.server, 40);
  WiFiManagerParameter mqtt_password("password", "mqtt password", config.password, 40);
  WiFiManagerParameter mqtt_user("password", "mqtt user", config.password, 40);
  WiFiManagerParameter mqtt_port("port", "mqtt port", port, 6);
  manager.setSaveConfigCallback(onSaveCallback);
  manager.setConnectTimeout(60);
  manager.addParameter(&mqtt_server);
  manager.addParameter(&mqtt_port);
  manager.addParameter(&mqtt_user);
  manager.addParameter(&mqtt_password);
  auto result = manager.autoConnect("ConfigAP");


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
    "192.168.100.10",
    config.user,
    config.password,
    CHIP_ID,
    1883
  );
  client->setKeepAlive(15);

  client->enableDebuggingMessages();
  client->enableLastWillMessage("doorlock/" CHIP_ID "/opener/status", "false", true);
  client->setMqttReconnectionAttemptDelay(5000);
  client->setWifiReconnectionAttemptDelay(60000);

}

void onOpenDoor(String data) {
  if(data == "true") {
    Serial.write(relON, sizeof(relON));
    delay(3000);
    Serial.write(relOFF, sizeof(relOFF));
  }
}

void onConnectionEstablished() {
  client->publish("doorlock/" CHIP_ID "/opener/status", "true", true);
  client->subscribe("doorlock/" CHIP_ID "/open", onOpenDoor);
}

void loop(void) {
  client->loop();
  delay(100);
}
