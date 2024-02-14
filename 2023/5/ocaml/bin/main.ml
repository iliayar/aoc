open Base
open Core

let split_list ~pred init =
  let rec impl ?(cur = []) ?(res = []) lst =
    match lst with
    | [] -> ( match cur with [] -> res | _ -> List.rev cur :: res)
    | e :: rest ->
        if pred e then
          match cur with
          | [] -> impl rest ~cur:[] ~res
          | _ -> impl rest ~cur:[] ~res:(List.rev cur :: res)
        else impl rest ~cur:(e :: cur) ~res
  in
  List.rev @@ impl init

type range = { l : int; r : int } [@@deriving sexp]

type mapping_elem = {
  src_begin : int;
  src_end : int;
  dst_begin : int;
  dst_end : int;
}
[@@deriving sexp]

type mapping = { src : string; dst : string; elems : mapping_elem list }
[@@deriving sexp]

let read_input =
  let parse_seeds l =
    let nums_str = String.chop_prefix_exn ~prefix:"seeds: " l in
    let nums = List.map ~f:Int.of_string @@ String.split ~on:' ' nums_str in
    let lc = List.chunks_of ~length:2 nums in
    List.map
      ~f:(fun p ->
        match p with
        | [ s; l ] -> { l = s; r = s + l }
        | _ -> failwith "unreachable")
      lc
  in
  let parse_mapping l =
    let l = String.chop_suffix_exn ~suffix:" map:" l in
    let l = String.split ~on:'-' l in
    match l with [ f; _; t ] -> (f, t) | _ -> failwith "unreachable"
  in
  let parse_section ls =
    let (f, t), ls =
      match ls with
      | m :: ls -> (parse_mapping m, ls)
      | _ -> failwith "unreachable"
    in
    let parse_l l =
      match List.map ~f:Int.of_string @@ String.split ~on:' ' l with
      | [ b1; b2; l ] ->
          { src_begin = b2; dst_begin = b1; src_end = b2 + l; dst_end = b1 + l }
      | _ -> failwith "unreachable"
    in
    let elems =
      List.sort ~compare:(fun { src_begin = bl; _ } { src_begin = br } ->
          Int.compare bl br)
      @@ List.map ~f:parse_l ls
    in
    { src = f; dst = t; elems }
  in
  let lines = Stdio.In_channel.read_lines "input.txt" in
  let sections = split_list lines ~pred:String.is_empty in
  let seeds_str, mappings_str =
    match sections with
    | hd :: tl -> (List.hd_exn hd, tl)
    | _ -> failwith "unreachable"
  in
  let init = parse_seeds seeds_str in
  let mappings = List.map ~f:parse_section mappings_str in
  (init, mappings)

(* let apply_mapping m n = *)
(*   let valid_ms = *)
(*     List.hd *)
(*     @@ List.filter *)
(*          ~f:(fun { src_begin; src_end; _ } -> n >= src_begin && n < src_end) *)
(*          m.elems *)
(*   in *)
(*   match valid_ms with Some m -> n - m.src_begin + m.dst_begin | None -> n *)
(**)
(* let apply_mappings m n = *)
(*   List.fold ~init:n ~f:(fun acc m -> apply_mapping m acc) m *)

let apply_mapping m rs =
  let rec impl l r m =
    if l >= r then []
    else
      match m with
      | [] -> [ { l; r } ]
      | m :: rest ->
          let s1 =
            if l < m.src_begin then [ { l; r = Int.min r m.src_begin } ] else []
          in
          let s2 =
            if l < m.src_end && r > m.src_begin then
              [
                {
                  l = Int.max l m.src_begin - m.src_begin + m.dst_begin;
                  r = Int.min r m.src_end - m.src_begin + m.dst_begin;
                };
              ]
            else []
          in
          s1 @ s2 @ impl (Int.max m.src_end l) r rest
  in
  List.fold ~init:[]
    ~f:(fun acc r ->
      let res = impl r.l r.r m.elems in
      (* Stdio.printf *)
      (*   !"%{sexp:mapping_elem list} %{sexp:range} -> %{sexp:range list}" *)
      (*   m.elems r res; *)
      acc @ res)
    rs

let apply_mappings m r =
  List.fold ~init:[ r ] ~f:(fun acc m -> apply_mapping m acc) m

let () =
  let init, mappings = read_input in
  Stdio.print_endline @@ Sexp.to_string_hum @@ [%sexp (init : range list)];
  Stdio.print_endline @@ Sexp.to_string_hum @@ [%sexp (mappings : mapping list)];
  let reses =
    List.fold ~init:[] ~f:(fun acc rs -> acc @ apply_mappings mappings rs) init
  in
  Stdio.print_endline @@ Sexp.to_string_hum @@ [%sexp (reses : range list)];
  let mr =
    List.min_elt ~compare:Int.compare @@ List.map ~f:(fun { l; _ } -> l) reses
  in
  Stdio.print_endline @@ Int.to_string @@ Option.value_exn mr
