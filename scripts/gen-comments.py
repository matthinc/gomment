import math
import sqlite3
import argparse

from idx_hierarchy import idx_to_hierarchy, hierarchy_to_idx, hierarchy_to_touched_at

parser = argparse.ArgumentParser(description='Bulk generate comments for testing')
parser.add_argument('file', metavar="FILE", type=str, help='target sqlite db file')
parser.add_argument('--thread', nargs=1, required=True, type=int, help='thread id to use for inserting comments')
parser.add_argument('--count', nargs=1, required=True, type=str, help='how many comments to insert at each layer per parent')
parser.add_argument('--purge', default=False, action=argparse.BooleanOptionalAction, help='whether to purge the thread before inserting')

args = parser.parse_args()
print(args)

counts = list(map(lambda e: int(e), args.count[0].split('-')))
thread_id = args.thread[0]

con = sqlite3.connect(args.file)
cur = con.cursor()

# start a transaction
cur.execute("BEGIN")
cur.execute("SELECT max(`comment_id`) FROM `comment`")
first_id = cur.fetchone()[0] + 1

print(f"starting with comment_id {first_id}")

if args.purge:
    cur.execute("DELETE FROM `comment` WHERE `thread_id` = ?", [thread_id])

comments = []

num_total = hierarchy_to_idx(counts, list(map(lambda e: e - 1, counts))) + 1
for relative_idx in range(num_total):
    md_idx = idx_to_hierarchy(counts, relative_idx)

    parent_id = None
    if len(md_idx) > 1:
        parent_id = first_id + hierarchy_to_idx(counts, md_idx[:len(md_idx)-1])

    num_children = 0
    if len(counts[len(md_idx):]) > 0:
        num_children = math.prod(counts[len(md_idx):])
    depth_level = len(md_idx) - 1
    created_at = first_id + relative_idx
    touched_at = first_id + hierarchy_to_touched_at(counts, md_idx)

    comments.append(
        (first_id + relative_idx, thread_id, parent_id, num_children, depth_level, created_at, touched_at, str(md_idx), str(md_idx))
    )


cur.executemany(
    "INSERT INTO `comment`(`comment_id`, `thread_id`, `parent_id`, `num_children`, `depth_level`, `created_at`, `touched_at`, `author`, `text`) values (?, ?, ?, ?, ?, ?, ?, ?, ?)",
    comments
)

if args.purge:
    cur.execute("UPDATE `thread` SET `num_total` = ?, `num_root` = ? WHERE `thread_id` = ?", [num_total, counts[0], thread_id])
else:
    cur.execute("UPDATE `thread` SET `num_total` = `num_total` + ?, `num_root` = `num_root` + ? WHERE `thread_id` = ?", [num_total, counts[0], thread_id])

cur.execute("COMMIT")

print(f"inserted {num_total} comments ({counts[0]} root)")
