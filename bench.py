import json
import random
import os
import subprocess
import sys
import time


TAGS = ['money', 'computer', 'shoe', 'dress', 'car', 'house', 'apple',
    'pencil', 'foot', 'star', 'planet', 'color', 'straw', 'battery',
    'controller', 'box', 'ram', 'harddrive', 'phone', 'toe', 'arm',
    'spanish', 'child', 'school', 'math', 'shorts', 'pizza', 'peroxide',
    'duster', 'light', 'table', 'up', 'cord', 'iron', 'bed', 'closet']
TAG_COUNT = (5, 16)
TAG_SCORE = (1, 10)
POOL_COUNT = 1000000


def make_item():
    data = {}

    for i in range(TAG_COUNT[0], TAG_COUNT[1]):
        tag = random.choice(TAGS)
        data[tag] = random.randint(TAG_SCORE[0], TAG_SCORE[1])

    return {
        'id': str(random.random()),
        'data': data,
    }


def run(source, pool, count=None):
    if count:
        pool = pool[:count]

    source = json.dumps(source)
    pool_len = len(pool)
    pool = json.dumps(pool)
    pool_file = '/tmp/gosinesim-{}.json'.format(pool_len)

    with open(pool_file, 'a') as j:
        j.write(pool)

    cmd = "./gosinesim --source='{}' --pool_file='{}' > /dev/null 2>&1".format(
        source, pool_file)

    print '*' * 80
    print "\nRunning {} tags\n".format(pool_len)

    for i in range(3):
        start = time.time()
        os.system(cmd)
        dif = time.time() - start

        print "run {} took {}".format(i + 1, dif)


if __name__ == '__main__':
    print "*" * 80
    print "\nBuilding the pool of {} tags".format(POOL_COUNT)

    start = time.time()
    source = make_item()
    pool = [make_item() for i in range(POOL_COUNT)]
    counts = [100, 1000, 10000, 50000, None]
    print "\ndone building the pool: {}".format(time.time() - start)

    [run(source, pool, c) for c in counts]
