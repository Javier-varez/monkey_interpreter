#ifndef _RUNTIME_CALLABLE_H
#define _RUNTIME_CALLABLE_H

#include <concepts>

namespace runtime {

template<typename T, typename R, typename ...Args>
concept Callable = requires(const T obj, Args... args) {
  { obj(args...) } -> std::same_as<R>;
};

}

#endif  // _RUNTIME_CALLABLE_H
