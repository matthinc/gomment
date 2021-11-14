import math

def idx_to_hierarchy(counts, index_in_level):
    for depth_level in range(len(counts)):
        level_total = math.prod(counts[:depth_level+1])

        if index_in_level >= level_total:
            index_in_level = index_in_level - level_total
            continue

        ret = []
        for lvl in reversed(range(depth_level + 1)):
            rows_in_level = counts[lvl]
            cur_idx = index_in_level % rows_in_level
            ret.insert(0, cur_idx)
            index_in_level = index_in_level // rows_in_level
        return ret

def hierarchy_to_idx(counts, md_idx):
    offset = 0
    for i in range(2,len(md_idx)+1):
        offset += math.prod(counts[:i-1])

    to_right = counts[:len(md_idx)]
    for i in range(len(to_right)):
        to_right[i] = math.prod(to_right[i+1:])

    relative_idx = 0
    for i in range(len(md_idx)):
        relative_idx = relative_idx + to_right[i] * md_idx[i]

    return offset + relative_idx

# get the touched_at date for a to-be generated parent
def hierarchy_to_touched_at(dimensions, hierarchy):
    # if the row is a leaf, there is no child to get the touched_at from
    if len(hierarchy) == len(dimensions):
        return hierarchy_to_idx(dimensions, hierarchy)

    # get the last possible index
    touched_at_h = [d - 1 for d in dimensions]

    # override with the start of the given hierarchy
    for idx, h in enumerate(hierarchy):
        touched_at_h[idx] = h

    return hierarchy_to_idx(dimensions, touched_at_h)
