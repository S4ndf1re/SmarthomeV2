/*
serialRelayTest.ino
*/

#include <ESP8266WiFi.h>
#include <PubSubClient.h>


void reconnect_wifi();

byte relON[] = {0xA0, 0x01, 0x01, 0xA2};
byte relOFF[] = {0xA0, 0x01, 0x00, 0xA1};


#define CHIP_ID "02"
#define MAX_TIMEOUT 20000

typedef struct {
  char server[40] = "";
  char user[40] = "";
  char password[40] = "";
  int port = 0;
} Config;

WiFiClient espClient;
PubSubClient client(espClient);

const String ssid = "<SSID>";
const String psk = "<PASSWORD>";
const char* mqtt_server = "<MQTT_SERVER_IP>";
const int mqtt_port = 1883;

Config config;




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


void onOpenDoor(char *topic, byte *payload, unsigned int length) {
  char buffer[length+1];
  copy_byte_to_cstr(payload, length, buffer, length+1);
  Serial.print("Received: ");
  Serial.println(buffer);
  if(strcmp(buffer,"true") == 0) {
    Serial.write(relON, sizeof(relON));
    delay(3000);
    Serial.write(relOFF, sizeof(relOFF));
  }
}

void onConnectionEstablished() {
  client.publish("doorlock/" CHIP_ID "/opener/status", "true", true);
  client.subscribe("doorlock/" CHIP_ID "/open");
}

void callback(char *topic, byte* payload, unsigned int length) {

  if (strcmp(topic, "doorlock/" CHIP_ID "/open") == 0) {
    onOpenDoor(topic, payload, length);
  }
}





void setup_wifi() {
  strcmp(config.server, mqtt_server);
  config.port = mqtt_port;
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

  WiFi.disconnect();
  
  // We start by connecting to a WiFi network
  Serial.println();
  Serial.print("Connecting to ");
  Serial.println(ssid);

  WiFi.mode(WIFI_STA);
  WiFi.begin(ssid, psk);

  int counter = 0;
  while (WiFi.status() != WL_CONNECTED) {
    if(counter > 60) {
      ESP.restart();
    }
    delay(1000);
    counter++;
    Serial.print(".");
  }

  randomSeed(micros());

  Serial.println("");
  Serial.println("WiFi connected");
  Serial.println("IP address: ");
  Serial.println(WiFi.localIP());
}


void setup(void) {
  Serial.begin(9600);
  Serial.print("Starting Connection");

  Serial.write(relOFF, sizeof(relOFF));
  delay(10);
  Serial.write(relOFF, sizeof(relOFF));


  setup_wifi();

  client.setServer(mqtt_server, mqtt_port);
  client.setCallback(callback);

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

void loop(void) {
  reconnect_wifi();

  if (!client.connected()) {
    reconnect_mqtt();
  }
  
  client.loop();
  delay(100);
}
