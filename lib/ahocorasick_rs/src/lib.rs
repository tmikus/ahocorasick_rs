extern crate aho_corasick;

use aho_corasick::{AhoCorasick, Match};
use libc::size_t;
use std::ffi::CStr;

#[no_mangle]
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
fn text_from_c(text: *const std::os::raw::c_char, text_len: usize) -> String {
    let rust_text = unsafe { std::slice::from_raw_parts(text as *const u8, text_len) };
    String::from_utf8_lossy(rust_text).to_string()
}

#[no_mangle]
pub extern "C" fn create_automaton(
    patterns: *const *const std::os::raw::c_char,
    num_patterns: usize,
) -> *mut AhoCorasick {
    let rust_patterns = (0..num_patterns)
        .map(|i| {
            let pattern = unsafe { CStr::from_ptr(*patterns.offset(i as isize)) };
            pattern.to_string_lossy().into_owned()
        })
        .collect::<Vec<String>>();

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
