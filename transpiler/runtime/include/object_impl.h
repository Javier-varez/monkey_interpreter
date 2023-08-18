#pragma once

#include <object.h>

namespace runtime {

template <typename... Args>
Object Object::operator()(const Args &...args) const noexcept {
  const Function f = std::get<Function>(val);
  return f(FnArgs{args...});
}

}
