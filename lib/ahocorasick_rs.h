#ifndef C_WRAPPER_H
#define C_WRAPPER_H

#include <stddef.h>

typedef struct AhoCorasickMatch {
    size_t end;
    size_t pattern_index;
    size_t start;
} AhoCorasickMatch;

typedef struct AhoCorasick AhoCorasick;

AhoCorasick* create_automaton(
    const char** patterns,
    size_t num_patterns
);

AhoCorasickMatch* find_iter(
    const AhoCorasick* automaton,
    const char* text,
    size_t text_len,
    long* found_count
);

void free_automaton(AhoCorasick* automaton);

#endif
