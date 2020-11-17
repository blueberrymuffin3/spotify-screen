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

#define DATA_BUFFER_COLORS 30
#define DATA_BUFFER_BYTES (DATA_BUFFER_COLORS * 2)
#define TIMEOUT 50

uint16_t readUInt16();

union _dataBuffer
{
  uint8_t bytes[DATA_BUFFER_COLORS * 2];
  uint16_t colors[DATA_BUFFER_COLORS];
} dataBuffer;

void setup()
{
  Serial.begin(4000000);
  Serial.setTimeout(TIMEOUT);

  tft.reset();
  tft.begin(tft.readID());
  tft.fillScreen(0xFFFF);
}

void loop()
{
  long start = millis();
  while (Serial.available() < 8)
    if (millis() - start >= TIMEOUT){
      for (uint8_t i = 0; i < 8; i++)
      {
        Serial.read();
      }      
      return;
    }

  uint16_t x1 = readUInt16();
  uint16_t y1 = readUInt16();
  uint16_t x2 = readUInt16();
  uint16_t y2 = readUInt16();

  tft.setAddrWindow(x1, y1, x2 - 1, y2 - 1);

  long bytesLeft = (long)(x2 - x1) * (long)(y2 - y1) * 2; // Must be able to go negative
  bool isFirst = true;

  if (bytesLeft > 153600) // Something is corrupt...
    return;

  while (bytesLeft > 0)
  {
    long maxLen = bytesLeft;
    if (maxLen > DATA_BUFFER_BYTES)
      maxLen = DATA_BUFFER_BYTES;

    size_t len = Serial.readBytes(dataBuffer.bytes, DATA_BUFFER_BYTES);
    // TODO: Should I check it this is odd?

    if (len == 0)
      return; // Timeout

    tft.pushColors(dataBuffer.colors, len / 2, isFirst);

    bytesLeft -= len;
    isFirst = false;
  }
}

uint16_t readUInt16()
{
  int low = Serial.read();
  int high = Serial.read();

  return (uint16_t)low | (uint16_t)high << 8;
}
