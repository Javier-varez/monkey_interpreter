cmake_minimum_required(VERSION 3.26)

project(monkey_transpiler)

set(CMAKE_CXX_STANDARD 20)
set(CMAKE_CXX_STANDARD_REQUIRED ON)
set(CMAKE_EXPORT_COMPILE_COMMANDS ON)

add_library(runtime STATIC
    src/array.cpp
    src/fn_args.cpp
    src/function.cpp
    src/object.cpp
    src/var_args.cpp
    src/hash_map.cpp)

target_include_directories(runtime PUBLIC include)

add_executable(main main.cpp)
target_link_libraries(main runtime)

if (BUILD_TESTS)

include(FetchContent)
FetchContent_Declare(
  googletest
  URL https://github.com/google/googletest/archive/03597a01ee50ed33e9dfd640b249b4be3799d395.zip
)
FetchContent_MakeAvailable(googletest)

enable_testing()

add_executable(
  runtime_test
  test/box_test.cpp
  test/rc_test.cpp
  test/vec_test.cpp
)

target_link_libraries(
  runtime_test
  runtime
  GTest::gtest_main
)

include(GoogleTest)
gtest_discover_tests(runtime_test)

endif()
