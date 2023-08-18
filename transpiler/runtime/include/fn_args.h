#pragma once

#include <object.h>
#include <var_args.h>

namespace runtime {

class FnArgs {
public:
  template <std::same_as<Object>... Args>
  explicit FnArgs(const Args &...args) noexcept;

  size_t len() const noexcept;

  Object operator[](size_t idx) const noexcept;

  using Iter = Iterator<const Object>;

  Iter begin() const noexcept;
  Iter end() const noexcept;

private:
  Vec<Object> args;
};

namespace detail {

template <typename... Args> size_t countArgs(const Args &...args) {
  const auto handleArg = [](const Object &arg) noexcept -> size_t {
    if (arg.type == ObjectType::VARARGS) {
      return arg.getVarArgs().len();
    }
    return 1;
  };

  return (handleArg(args) + ...);
}

} // namespace detail

template <std::same_as<Object>... Args>
FnArgs::FnArgs(const Args &...args) noexcept
    : args(
          [args...](auto pusher) noexcept -> void {
            const auto handleArg = [pusher]<typename T>(const T arg) {
              static_assert(std::same_as<T, Object>, "Invalid arg type");
              if (arg.type == ObjectType::VARARGS) {
                // Varargs are unwrapped here in the call site
                for (const Object &inner : arg.getVarArgs()) {
                  pusher.push(inner);
                }
              } else {
                pusher.push(arg);
              }
            };

            (handleArg(args), ...);
          },
          detail::countArgs(args...)) {}

} // namespace runtime
