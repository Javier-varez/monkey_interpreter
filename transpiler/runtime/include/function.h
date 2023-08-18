#pragma once

#include <box.h>
#include <rc.h>

namespace runtime {

struct Object;
class FnArgs;

template <typename T, T> struct ConstexprLit {};

class Function final {
private:
  struct Callable {
    virtual Object call(const FnArgs &args) const noexcept = 0;
    virtual ~Callable() noexcept = default;
  };

  template <typename T, size_t NumArgs, bool HasVarArgs>
  struct CallableImpl final : public Callable {
    CallableImpl(T callable) : callable{callable} {}

    Object call(const FnArgs &args) const noexcept final;

    T callable;
  };

public:
  template <typename T, size_t NumArgs, bool HasVarArgs>
  Function(ConstexprLit<size_t, NumArgs>, ConstexprLit<bool, HasVarArgs>,
           T &&callable)
      : callable{Marker<CallableImpl<T, NumArgs, HasVarArgs>>{}, callable} {}

  Object operator()(const FnArgs &args) const noexcept;

private:
  Rc<Callable> callable;
};

} // namespace runtime
