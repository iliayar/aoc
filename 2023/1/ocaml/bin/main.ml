open Batteries

let read_input : string list = List.of_enum (File.lines_of "input.txt")

let sem_num_mapping =
  [
    ("one", 1);
    ("two", 2);
    ("three", 3);
    ("four", 4);
    ("five", 5);
    ("six", 6);
    ("seven", 7);
    ("eight", 8);
    ("nine", 9);
    ("zero", 0);
  ]

let try_match_sem_num s =
  match
    List.find_opt (fun (sem, _) -> String.starts_with s sem) sem_num_mapping
  with
  | Some (_, n) -> Some n
  | None -> None

let try_get_num s =
  match Seq.first (String.to_seq s) with
  | '0' .. '9' as ch -> Some (Char.code ch - Char.code '0')
  | _ -> try_match_sem_num s

let rec extract_nums s =
  let len = String.length s in
  if len == 0 then []
  else
    let rest = extract_nums (String.sub s 1 (len - 1)) in
    match try_get_num s with Some n -> n :: rest | None -> rest

let rec last = function
  | [] -> failwith "list is empty"
  | e :: [] -> e
  | _ :: es -> last es

let first = function [] -> failwith "list is empty" | e :: _ -> e

let calc ls =
  let l = last ls in
  let f = first ls in
  (f * 10) + l

let sum ls = List.fold_left ( + ) 0 ls

let () =
  let input = read_input in
  let nums = List.map extract_nums input in
  print_endline (Int.to_string (sum (List.map calc nums)))
