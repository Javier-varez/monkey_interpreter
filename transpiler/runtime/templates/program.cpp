#include <runtime.h>

using std::literals::operator""sv;
using std::literals::operator""s;

int main() {
  {{range .Statements}} {{Transpile .}} {{end}}
  return 0;
}
