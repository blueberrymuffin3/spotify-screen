#include <Arduino.h>
#include <SPI.h>
#include <Adafruit_I2CDevice.h>
#include <Adafruit_GFX.h>
#include <Adafruit_TFTLCD.h>
#include <SPI.h>

#define LCD_CS A3
#define LCD_CD A2
#define LCD_WR A1
#define LCD_RD A0

Adafruit_TFTLCD tft(LCD_CS, LCD_CD, LCD_WR, LCD_RD, A4);

#define DATA_BUFFER_COLORS 127
#define DATA_BUFFER_BYTES (DATA_BUFFER_COLORS * 2)
#define TIMEOUT 250

uint16_t readUInt16();

union _dataBuffer
{
  uint8_t bytes[DATA_BUFFER_COLORS * 2];
  uint16_t colors[DATA_BUFFER_COLORS];
} dataBuffer;

void setup()
{
  Serial.begin(1000000);
  Serial.setTimeout(TIMEOUT);

  tft.reset();

  uint16_t id = tft.readID();
  
  tft.begin(id);
  tft.fillScreen(0x8888);
}

void loop()
{
  while (Serial.available() < 8); // Wait for a header

  int32_t x1 = readUInt16();
  int32_t y1 = readUInt16();
  int32_t x2 = readUInt16();
  int32_t y2 = readUInt16();

  if (x1 >= x2 || y1 >= y2 || x1 < 0 || y1 < 0 || x2 > 240 || y2 > 320) // Something is corrupt...
  {
    Serial.write('C');
    delay(2000);
    while (Serial.available())
    {
      Serial.read();
    }

    return;
  }

  tft.setAddrWindow(x1, y1, x2 - 1, y2 - 1);

  int32_t bytesLeft = (int32_t)(x2 - x1) * (int32_t)(y2 - y1) * 2; // Must be able to go negative
  bool isFirst = true;

  while (bytesLeft > 0)
  {
    int32_t maxLen = bytesLeft;
    if (maxLen > DATA_BUFFER_BYTES)
      maxLen = DATA_BUFFER_BYTES;

    int32_t len = Serial.readBytes(dataBuffer.bytes, DATA_BUFFER_BYTES);

    if (len == 0) // Timeout
    {
      Serial.write('T');
      return;
    }

    tft.pushColors(dataBuffer.colors, len / 2, isFirst);

    bytesLeft -= len;
    isFirst = false;
  }

  Serial.write('O');
}

uint16_t readUInt16()
{
  int low = Serial.read();
  int high = Serial.read();

  return (uint16_t)low | (uint16_t)high << 8;
}
