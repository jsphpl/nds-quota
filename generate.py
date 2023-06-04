#!/usr/bin/env python

import random
import argparse
import os
import json
import string

def mkdir(path: str) -> None:
    try:
        os.mkdir(path)
    except FileExistsError:
        pass

def random_chars(length: int) -> str:
    letters = string.ascii_uppercase
    return "".join(random.choice(letters) for i in range(length))

def random_numbers(length: int) -> str:
    letters = string.digits
    return "".join(random.choice(letters) for i in range(length))

def generate(path: str, quota_kib: int) -> None:
    id = '{}-{}'.format(random_chars(4), random_numbers(4))
    value = {
        "id": id,
        "quota_kib": quota_kib,
        "used_kib": 0,
        "devices": [],
    }

    with open(os.path.join(path, id), "w") as f:
        json.dump(value, f)

    print(id)

if __name__ == "__main__":
    parser = argparse.ArgumentParser()
    parser.add_argument("dir", help="Directory path to put the accounts in (will be created)", type=str)
    parser.add_argument("--count", help="Number of accounts to generate", type=int, default=1)
    parser.add_argument("--quota", help="Quota per account in KiB", type=int, default=102400)
    args = parser.parse_args()

    mkdir(args.dir)
    for i in range(args.count):
        generate(args.dir, args.quota)

