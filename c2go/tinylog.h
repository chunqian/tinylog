/**---------------------------------------------------------
 * name: tinylog.h
 * author: shenchunqian
 * created: 2022-08-19
 ---------------------------------------------------------*/

#ifndef TINYLOG_H
#define TINYLOG_H

#ifndef SUPPORT_TINY_LOG
  #define log_debug(...)
  #define log_info(...)
  #define log_warn(...)
  #define log_error(...)
  #define log_fatal(...)
  #define log_message(...)
#else
  int printf(const char *__restrict, ...);
  #define log_debug(...) printf(0, "DEBUG", __VA_ARGS__)
  #define log_info(...) printf(0, "INFO", __VA_ARGS__)
  #define log_warn(...) printf(0, "WARN", __VA_ARGS__)
  #define log_error(...) printf(0, "ERROR", __VA_ARGS__)
  #define log_fatal(...) printf(0, "FATAL", __VA_ARGS__)
  #define log_message(...) printf(0, "MESSAGE", __VA_ARGS__)
#endif

#endif // TINYLOG_H
