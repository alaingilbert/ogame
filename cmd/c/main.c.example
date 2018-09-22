#include <stdio.h>
#include <stdlib.h>
#include "ogame.h"
int main() {
  char* universe = getenv("UNIVERSE");
  char* username = getenv("USERNAME");
  char* password = getenv("PASSWORD");
  char* error_msg;
  error_msg = OGame(universe, username, password);
  if (error_msg) {
    printf("Error: %s\n", error_msg);
    exit(1);
  }

  struct GetPlanet_return p = GetPlanet(123);
  if (p.r11) {
    printf("Error: %s\n", p.r11);
  }

  int is_under_attack = IsUnderAttack();
  if (is_under_attack == 1) {
    printf("Attack detected\n");
  } else {
    printf("No attack detected\n");
  }
}
