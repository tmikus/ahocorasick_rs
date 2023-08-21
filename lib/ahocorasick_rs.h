#ifndef C_WRAPPER_H
#define C_WRAPPER_H

#include <stddef.h>

typedef struct AhoCorasick AhoCorasick;

AhoCorasick* create_automaton(const char** patterns, size_t num_patterns);
void free_automaton(AhoCorasick* automaton);
size_t search_automaton(const AhoCorasick* automaton, const char* text, size_t text_len, size_t* matches);

#endif
