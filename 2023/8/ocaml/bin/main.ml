open Base
open Core

type direction = L | R [@@deriving sexp]

let direction_from_char c =
  match c with 'L' -> L | 'R' -> R | _ -> failwith "unreachable"

let directions_from_string s =
  Array.map ~f:direction_from_char @@ String.to_array s

let kBegin = "AAA"
let kEnd = "ZZZ"

let read_input =
  let lines = Stdio.In_channel.read_lines "input.txt" in
  let parse_edge s =
    match String.split ~on:' ' s with
    | [ from; _; ls; rs ] ->
        let l =
          String.chop_prefix_exn ~prefix:"("
          @@ String.chop_suffix_exn ~suffix:","
          @@ ls
        in
        let r = String.chop_suffix_exn ~suffix:")" rs in
        (from, (l, r))
    | _ -> failwith "unreachable"
  in
  let path, edges =
    match lines with
    | path :: _ :: edges ->
        (directions_from_string path, List.map ~f:parse_edge edges)
    | _ -> failwith "unreachable"
  in
  (path, edges)

let rec gcd a = function 0 -> a | b -> gcd b (a % b)
let lcm x y = x * y / gcd x y

let walk path edges init =
  let len = List.length init in
  let c = ref len in
  let periods = Array.create ~len 0 in
  let was = Array.create ~len false in
  let rec impl es i =
    (* Stdio.printf !"%{sexp:string list}" es; *)
    List.iteri
      ~f:(fun j e ->
        if String.is_suffix ~suffix:"Z" e && not was.(j) then (
          Stdio.printf "%d: %d\n" j i;
          c := !c - 1;
          periods.(j) <- i;
          was.(j) <- true))
      es;
    if !c = 0 then ()
    else
      let dir = path.(i % Array.length path) in
      let next e =
        let left, right = Option.value_exn @@ Map.find edges e in
        match dir with L -> left | R -> right
      in
      impl (List.map ~f:next es) (i + 1)
  in
  impl init 0;
  List.of_array periods

let () =
  let path, edges_lst = read_input in
  let edges =
    List.fold
      ~init:(Map.empty (module String))
      ~f:(fun acc (from, lr) -> Map.add_exn acc ~key:from ~data:lr)
      edges_lst
  in
  let init_es =
    List.filter_map
      ~f:(fun (e, _) -> if String.is_suffix ~suffix:"A" e then Some e else None)
      edges_lst
  in
  Stdio.print_endline @@ Sexp.to_string_hum
  @@ Map.sexp_of_m__t (module String) [%sexp_of: string * string] edges;
  let periods = walk path edges init_es in
  let res = List.fold ~init:1 ~f:(fun acc x -> lcm acc x) periods in
  Stdio.printf !"%{sexp:int list}\n" periods;
  Stdio.print_endline @@ Int.to_string res
