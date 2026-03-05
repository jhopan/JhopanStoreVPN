cmake_minimum_required(VERSION 3.18)
project(spawn_helper C)

add_library(spawn_helper SHARED spawn_helper.c)

find_library(log-lib log)
target_link_libraries(spawn_helper ${log-lib})
