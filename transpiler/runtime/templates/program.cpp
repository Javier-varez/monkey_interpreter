#include <runtime.h>

namespace runtime {

using std::literals::operator""sv;
using std::literals::operator""s;

void _mainProgram() {
  {{range .Statements}}
    {{ Transpile . }}
  {{end}}
}

}

int main() {
  runtime::_mainProgram();
  return 0;
}
