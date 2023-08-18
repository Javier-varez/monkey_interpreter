#include <function.h>

#include <object.h>
#include <var_args.h>
#include <fn_args.h>

namespace runtime {

Object Function::operator()(const FnArgs &args) const noexcept {
  const auto result = callable->call(args);
  return result;
}

}
