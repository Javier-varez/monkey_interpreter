#include <fn_args.h>

namespace runtime {

size_t FnArgs::len() const noexcept { return args.size(); }

Object FnArgs::operator[](size_t idx) const noexcept {
  using std::literals::operator""sv;
  check(idx < args.size(), "Out of bounds index to FnArgs object"sv);
  return args[idx];
}

FnArgs::Iter FnArgs::begin() const noexcept { return args.begin(); }
FnArgs::Iter FnArgs::end() const noexcept { return args.end(); }

}
