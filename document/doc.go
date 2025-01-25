// Package document provides document types and processing functionality.
//
// Processing steps:
//  1. Collect full paths of all (referenceable) members (1st recursion).
//  2. Extract doc tests (2nd recursion).
//  3. Filter and restructure according to re-exports
//     a) Collect re-exports per package
//     b) Traverse the tree and build up the re-structured package(s).
//     Also collects mapping from old to new paths (3rd recursion).
//  4. Traverse the new tree to collect all link paths (4th recursion).
//  5. Traverse the old tree and replace refs by full path placeholders.
//     Also changes the doc-strings in the new tree, as the new tree contains pointers into the old tree (5th recursion).
//  6. Render the new tree, replacing full path placeholders by actual relative links (6th recursion).
package document
