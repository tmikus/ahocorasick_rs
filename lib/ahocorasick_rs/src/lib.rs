extern crate aho_corasick;

use aho_corasick::{AhoCorasick, AhoCorasickKind, Match, MatchKind, StartKind};
use libc::size_t;
use std::ffi::CStr;

#[repr(C)]
pub struct AhoCorasickMatch {
    end: size_t,
    pattern_index: size_t,
    start: size_t,
}

#[inline]
fn aho_corasick_match_from_match(m: Match) -> AhoCorasickMatch {
    AhoCorasickMatch {
        end: m.end() as size_t,
        pattern_index: m.pattern().as_usize() as size_t,
        start: m.start() as size_t,
    }
}

#[inline]
fn patterns_from_c(
    patterns: *const *const std::os::raw::c_char,
    num_patterns: usize,
) -> Vec<String> {
    (0..num_patterns)
        .map(|i| {
            let pattern = unsafe { CStr::from_ptr(*patterns.offset(i as isize)) };
            pattern.to_string_lossy().into_owned()
        })
        .collect::<Vec<_>>()
}

#[inline]
fn text_from_c(text: *const std::os::raw::c_char, text_len: usize) -> &'static [u8] {
    unsafe { std::slice::from_raw_parts(text as *const u8, text_len) }
}

#[repr(C)]
pub struct AhoCorasickBuilderOptions {
    ascii_case_insensitive: i32,
    byte_classes: i32,
    dense_depth: *const usize,
    kind: *const usize,
    match_kind: usize,
    prefilter: i32,
    start_kind: usize,
}

impl AhoCorasickBuilderOptions {
    fn get_kind(&self) -> Option<AhoCorasickKind> {
        if self.kind.is_null() {
            return None;
        }
        match unsafe { *self.kind } {
            1 => Some(AhoCorasickKind::NoncontiguousNFA),
            2 => Some(AhoCorasickKind::ContiguousNFA),
            3 => Some(AhoCorasickKind::DFA),
            _ => None,
        }
    }

    fn get_match_kind(&self) -> MatchKind {
        match self.match_kind {
            1 => MatchKind::Standard,
            2 => MatchKind::LeftmostLongest,
            3 => MatchKind::LeftmostFirst,
            _ => panic!("Invalid match kind"),
        }
    }

    fn get_start_kind(&self) -> StartKind {
        match self.start_kind {
            1 => StartKind::Both,
            2 => StartKind::Unanchored,
            3 => StartKind::Anchored,
            _ => panic!("Invalid start kind"),
        }
    }
}

#[no_mangle]
pub extern "C" fn build_automaton(
    patterns: *const *const std::os::raw::c_char,
    num_patterns: usize,
    options: *const AhoCorasickBuilderOptions,
) -> *mut AhoCorasick {
    let rust_patterns = patterns_from_c(patterns, num_patterns);
    let rust_options = unsafe { &*options };
    let mut builder = AhoCorasick::builder();
    println!(
        "ascii_case_insensitive: {}",
        rust_options.ascii_case_insensitive
    );
    builder.ascii_case_insensitive(rust_options.ascii_case_insensitive != 0);
    builder.byte_classes(rust_options.byte_classes != 0);
    if !rust_options.dense_depth.is_null() {
        builder.dense_depth(unsafe { *rust_options.dense_depth });
    }
    builder.kind(rust_options.get_kind());
    builder.match_kind(rust_options.get_match_kind());
    builder.prefilter(rust_options.prefilter != 0);
    builder.start_kind(rust_options.get_start_kind());
    match builder.build(&rust_patterns) {
        Ok(automaton) => Box::into_raw(Box::new(automaton)),
        Err(_) => std::ptr::null_mut(),
    }
}

#[no_mangle]
pub extern "C" fn create_automaton(
    patterns: *const *const std::os::raw::c_char,
    num_patterns: usize,
) -> *mut AhoCorasick {
    let rust_patterns = patterns_from_c(patterns, num_patterns);
    match AhoCorasick::new(&rust_patterns) {
        Ok(automaton) => Box::into_raw(Box::new(automaton)),
        Err(_) => std::ptr::null_mut(),
    }
}

#[no_mangle]
pub extern "C" fn free_automaton(automaton: *mut AhoCorasick) {
    if automaton.is_null() {
        return;
    }
    unsafe {
        let _ = Box::from_raw(automaton);
    }
}

#[no_mangle]
pub extern "C" fn find(
    automaton: *const AhoCorasick,
    text: *const std::os::raw::c_char,
    text_len: usize,
) -> *mut AhoCorasickMatch {
    let rust_text = text_from_c(text, text_len);
    let automaton_ref = unsafe { &*automaton };
    let result = automaton_ref
        .find(&rust_text)
        .map(aho_corasick_match_from_match);
    match result {
        Some(m) => Box::into_raw(Box::new(m)),
        None => std::ptr::null_mut(),
    }
}

#[no_mangle]
pub extern "C" fn find_iter(
    automaton: *const AhoCorasick,
    text: *const std::os::raw::c_char,
    text_len: usize,
    found_count: *mut size_t,
) -> *mut AhoCorasickMatch {
    let rust_text = text_from_c(text, text_len);
    let automaton_ref = unsafe { &*automaton };
    let mut result = automaton_ref
        .find_iter(&rust_text)
        .map(aho_corasick_match_from_match)
        .collect::<Vec<_>>();
    unsafe {
        *found_count = result.len();
    }
    let ptr = result.as_mut_ptr();
    std::mem::forget(result);
    ptr
}

#[no_mangle]
pub extern "C" fn is_match(
    automaton: *const AhoCorasick,
    text: *const std::os::raw::c_char,
    text_len: usize,
) -> bool {
    let rust_text = text_from_c(text, text_len);
    let automaton_ref = unsafe { &*automaton };
    automaton_ref.is_match(&rust_text)
}
