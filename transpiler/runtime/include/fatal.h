#pragma once

#include <utility>
#include <iostream>

namespace runtime {

template <typename... Args>
[[noreturn]] void fatal(Args &&...args) noexcept {
  const auto print = []<typename T>(T &&arg) {
    std::cout << std::forward<T>(arg);
    return true;
  };

  using std::literals::operator""sv;
  std::cout << "Failed assertion: "sv;

  (print(std::forward<Args>(args)) && ...);
  std::exit(-1);
}

template <typename... Args>
constexpr void check(const bool condition, Args &&...args) noexcept {
  if (!condition) {
    fatal(std::forward<Args>(args)...);
  }
}

}
