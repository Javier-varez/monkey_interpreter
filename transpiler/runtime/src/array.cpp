#include <array.h>

#include <callable.h>
#include <object.h>
#include <var_args.h>

namespace runtime {

Object Array::operator[](size_t index) const noexcept {
  check(index < len(), "Out of bounds access to array.");
  return data[index];
}

size_t Array::len() const noexcept { return data.size(); }

Array::Iter Array::begin() const noexcept { return data.begin(); }

Array::Iter Array::end() const noexcept { return data.end(); }

Array Array::makeFromRange(int64_t start, int64_t end) noexcept {
  const auto abs = [](auto arg) {
    if (arg < 0)
      return -arg;
    return arg;
  };
  const size_t sizeHint = abs(end - start);

  return Array{[start, end](LargeVec<Object>::Pusher pusher) noexcept -> void {
                 int64_t current = start;
                 if (start > end) {
                   while (current > end) {
                     pusher.push(Object::makeInt(current--));
                   }
                 } else {
                   while (current < end) {
                     pusher.push(Object::makeInt(current++));
                   }
                 }
               },
               sizeHint};
}

Array Array::makeFromIters(const Iter begin, const Iter end) noexcept {
  return Array{[begin, end](LargeVec<Object>::Pusher pusher) noexcept -> void {
                 auto next = begin;
                 while (next < end) {
                   pusher.push(*next);
                   next++;
                 }
               },
               static_cast<size_t>(end - begin)};
}

Array Array::push(const Object &newObj) const noexcept {
  return Array{data.copyAppend(newObj)};
}

}
