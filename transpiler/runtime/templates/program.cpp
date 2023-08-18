#include <runtime.h>

namespace runtime {

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
