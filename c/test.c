/**---------------------------------------------------------
 * name: test.c
 * author: shenchunqian
 * created: 2022-08-19
 ---------------------------------------------------------*/

#define SUPPORT_TINY_LOG

#include "tinylog.h"

int main()
{
  /* code */
  char *msg = "Hello C!";
  log_message("msg: {}", msg);
  return 0;
}
