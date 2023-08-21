extern crate aho_corasick;

use aho_corasick::AhoCorasick;
use std::ffi::CStr;

#[no_mangle]
pub extern "C" fn create_automaton(patterns: *const *const std::os::raw::c_char, num_patterns: usize) -> *mut AhoCorasick {
    let rust_patterns = (0..num_patterns)
        .map(|i| {
            let pattern = unsafe { CStr::from_ptr(*patterns.offset(i as isize)) };
            pattern.to_string_lossy().into_owned()
        })
        .collect::<Vec<String>>();

    match AhoCorasick::new(&rust_patterns) {
        Ok(automaton) => Box::into_raw(Box::new(automaton)),
        Err(_) => {
            std::ptr::null_mut()
        },
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
pub extern "C" fn search_automaton(
    automaton: *const AhoCorasick,
    text: *const std::os::raw::c_char,
    text_len: usize,
    matches: *mut usize,
) -> usize {
    let rust_text = unsafe { std::slice::from_raw_parts(text as *const u8, text_len) };
    let rust_text = String::from_utf8_lossy(rust_text).to_string();

    let automaton_ref = unsafe { &*automaton };
    let result = automaton_ref
        .find_iter(&rust_text)
        .map(|m| m.pattern().as_usize())
        .collect::<Vec<_>>();
    result.

    for (i, match_index) in result.iter().enumerate() {
        unsafe {
            *matches.offset(i as isize) = match_index.as_usize();
        }
    }

    result.len()
}