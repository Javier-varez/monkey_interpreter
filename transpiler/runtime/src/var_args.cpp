#include <var_args.h>

namespace runtime {

VarArgs::VarArgs(const Iter begin, const Iter end) noexcept
    : args([begin, end](auto pusher) noexcept -> void {
        Iter next = begin;
        while (next < end) {
          pusher.push(*next);
          next++;
        }
      }) {}

size_t VarArgs::len() const noexcept { return args.size(); }

Object VarArgs::operator[](size_t idx) const noexcept { return args[idx]; }

VarArgs::Iter VarArgs::begin() const noexcept { return args.begin(); }

VarArgs::Iter VarArgs::end() const noexcept { return args.end(); }

} // namespace runtime
