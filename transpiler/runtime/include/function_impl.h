#pragma once

#include <fn_args.h>
#include <function.h>
#include <object.h>
#include <var_args.h>

namespace runtime {

template <size_t NumArgs, typename C, typename... Args>
auto expandAndCall(const Iterator<const Object> argIter, C &&callable,
                   Args... args) {
  if constexpr (NumArgs == 0) {
    return callable(std::forward<Args>(args)...);
  } else {
    const auto nextIter = argIter + 1;
    return expandAndCall<NumArgs - 1>(nextIter, std::forward<C>(callable),
                                      args..., *argIter);
  }
}

template <size_t NumArgs, typename C, typename... Args>
auto expandAndCallWithVarArgs(const Iterator<const Object> argIter,
                              const Iterator<const Object> argIterEnd,
                              C &&callable, Args... args) {
  if constexpr (NumArgs == 0) {
    const Object varArgs{Object::makeVarargs(VarArgs{argIter, argIterEnd})};
    return callable(std::forward<Args>(args)..., varArgs);
  } else {
    const auto nextIter = argIter + 1;
    return expandAndCallWithVarArgs<NumArgs - 1>(
        nextIter, argIterEnd, std::forward<C>(callable), args..., *argIter);
  }
}

template <typename T, size_t NumArgs, bool HasVarArgs>
Object Function::CallableImpl<T, NumArgs, HasVarArgs>::call(
    const FnArgs &args) const noexcept {
  using std::literals::operator""sv;
  if constexpr (!HasVarArgs) {
    check(args.len() == NumArgs, "Callable takes "sv, NumArgs,
          " arguments, but "sv, args.len(), " were given");
    return expandAndCall<NumArgs>(args.begin(), callable);
  } else {
    check(args.len() >= NumArgs, "Callable takes at least "sv, NumArgs,
          " arguments, but only "sv, args.len(), " were given");
    return expandAndCallWithVarArgs<NumArgs>(args.begin(), args.end(),
                                             callable);
  }
}

} // namespace runtime
