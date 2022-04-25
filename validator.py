import os
import subprocess
import sys

import difflib

from tqdm import tqdm

def measure_time(num_of_iters: int, program_file: str,
                 config_file: str, res_a_path: str, res_n_path: str) -> int:
    cur_a, prev_a = None, None
    cur_n, prev_n = None, None
    time_taken = [sys.maxsize for _ in range(4)]

    for _ in tqdm(range(num_of_iters)):
        outputs = subprocess.run(
            [program_file, config_file],
            capture_output=True)

        try:
            out = outputs.stdout.decode("utf-8").split('\n')
            for idx, elm in enumerate(out):
                if not elm:
                    break
                out[idx] = float(elm.split(":")[-1].rstrip("ms"))

            indexers, mergers, writing, total, _ = out
        except Exception as e:
            print(e)
            return

        with open(res_a_path, "r") as f:
            res_a = f.readlines()
        prev_a, cur_a = cur_a, res_a

        with open(res_n_path, "r") as f:
            res_n = f.readlines()
        prev_n, cur_n = cur_n, res_n

        if prev_a is not None:
            if prev_a != cur_a or prev_n != cur_n:
                print("Result does not match")
                print("diff for the firts file\n")
                for line in difflib.unified_diff(prev_a, cur_a, fromfile='prev_n', tofile='cur_n'):
                    print(line)
                return -1
        
        time_taken = [
            min(time_taken[0], total),
            min(time_taken[1], indexers),
            min(time_taken[2], mergers),
            min(time_taken[3], writing),
        ]

    return time_taken

if __name__ == "__main__":

    args = sys.argv
    args_num = 2
    config_file = "index.cfg"
    program = "./main"
    res_a = "res_a.txt"
    res_n = "res_n.txt"

    if sys.platform == "win32":
        args_num = 3
        program = ".\\bin\\q.exe"

    if len(args) != args_num:
        print("usage: python3 validator.py number_of_iterations_to_measure")
        print("ATTENTION! If your operation system is Windows, please provide one more parameter - mingw64 folder path")
        exit(1)

    if sys.platform  == "win32":
        os.environ["PATH"] = sys.argv[2]+'\\bin'

    max_iters = int(sys.argv[1])

    res = measure_time(max_iters, program, config_file, res_a, res_n)
    if res != -1:
        print(res)