#ifndef C_WRAPPER_H
#define C_WRAPPER_H

#include <stddef.h>

typedef struct AhoCorasick AhoCorasick;

typedef struct AhoCorasickBuilderOptions {
    int ascii_case_insensitive;
    int byte_classes;
    size_t* dense_depth;
    size_t* kind;
    size_t match_kind;
    int prefilter;
    size_t start_kind;
} AhoCorasickBuilderOptions;

typedef struct AhoCorasickMatch {
    size_t end;
    size_t pattern_index;
    size_t start;
} AhoCorasickMatch;

AhoCorasick* build_automaton(
    const char** patterns,
    size_t num_patterns,
    const AhoCorasickBuilderOptions* builder
);

AhoCorasick* create_automaton(
    const char** patterns,
    size_t num_patterns
);

AhoCorasickMatch* find(
    const AhoCorasick* automaton,
    const char* text,
    size_t text_len
);

AhoCorasickMatch* find_iter(
    const AhoCorasick* automaton,
    const char* text,
    size_t text_len,
    long* found_count
);

void free_automaton(AhoCorasick* automaton);

int get_kind(const AhoCorasick* automaton);

int is_match(
    const AhoCorasick* automaton,
    const char* text,
    size_t text_len
);

#endif
