import subprocess
import sys
import os


def get_executable_path():
    exec_names = {
        'win32': 'modo_win.exe',
        'darwin': 'modo_macos',
        'linux': 'modo_linux',
    }
    platform = os.sys.platform
    exec_name = exec_names[platform]
    return os.path.join(os.path.dirname(__file__), "bin", exec_name)


def main():
    executable_path = get_executable_path()
    if os.path.exists(executable_path):
        subprocess.run([executable_path] + sys.argv[1:])
    else:
        print(f"Executable not found at {executable_path}")


if __name__ == "__main__":
    main()
