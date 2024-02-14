open Base

type direction = L | R | T | B [@@deriving ord, sexp]

let opposite = function L -> R | R -> L | T -> B | B -> T

module Cell = struct
  type t = Pipe of direction * direction | Ground | Start [@@deriving sexp]

  let of_char : char -> t = function
    | '|' -> Pipe (T, B)
    | '-' -> Pipe (L, R)
    | 'L' -> Pipe (T, R)
    | 'J' -> Pipe (T, L)
    | '7' -> Pipe (L, B)
    | 'F' -> Pipe (R, B)
    | '.' -> Ground
    | 'S' -> Start
    | _ -> failwith "unreachable"

  let is_pipe = function Pipe _ -> true | _ -> false
  let is_start = function Start -> true | _ -> false
  let is_ground = function Ground -> true | _ -> false

  let has_diretion dir = function
    | Pipe (f, s) -> compare_direction dir f = 0 || compare_direction dir s = 0
    | _ -> failwith "invalid argument"
end

let read_input =
  let lines = Stdio.In_channel.read_lines "input.txt" in
  let map =
    Array.of_list
    @@ List.map
         ~f:(fun l -> Array.map ~f:Cell.of_char @@ String.to_array l)
         lines
  in
  map

let check_connected map pos direction =
  let x, y = pos in
  let xn, yn =
    match direction with
    | L -> (x - 1, y)
    | R -> (x + 1, y)
    | T -> (x, y - 1)
    | B -> (x, y + 1)
  in
  let h = Array.length map in
  let w = Array.length map.(0) in
  if xn < 0 || xn >= w || yn < 0 || yn >= h then None
  else
    let ret b = if b then Some (xn, yn) else None in
    let c = map.(y).(x) in
    let cn = map.(yn).(xn) in
    Cell.(
      if is_ground c || is_ground cn then None
      else
        ret
          ((is_start c || has_diretion direction c)
          && (is_start cn || has_diretion (opposite direction) cn)))

let all_directions = [ L; R; T; B ]

let find_start map =
  let h = Array.length map in
  let w = Array.length map.(0) in
  let rec xloop x =
    let rec yloop y =
      if y = h then None
      else if Cell.is_start map.(y).(x) then Some (x, y)
      else yloop (y + 1)
    in
    if x = w then None
    else match yloop 0 with Some c -> Some c | None -> xloop (x + 1)
  in
  xloop 0

let map_array map v =
  Array.init (Array.length map) ~f:(fun _ ->
      Array.create ~len:(Array.length map.(0)) v)

let print_map to_string map =
  let char_map = Array.map ~f:(Array.map ~f:to_string) map in
  Array.iter
    ~f:(fun l ->
      let s = String.concat_array ~sep:" " l in
      Stdio.print_endline s)
    char_map

let bfs map s =
  let dists = map_array map 0 in
  let was = map_array map false in
  let impl queue =
    while not @@ Queue.is_empty queue do
      let x, y = Queue.dequeue_exn queue in
      let curd = dists.(y).(x) in
      (* Stdio.printf !"Visit: %{sexp:int * int}, D: %d\n" (x, y) curd; *)
      (* print_map Int.to_string dists; *)
      (* print_map Bool.to_string was; *)
      let ns =
        List.filter ~f:(fun (x, y) -> not was.(y).(x))
        @@ List.filter_map ~f:(check_connected map (x, y)) all_directions
      in
      List.iter
        ~f:(fun (xn, yn) ->
          dists.(yn).(xn) <- curd + 1;
          was.(y).(x) <- true;
          (* Stdio.printf !"Enqueue: %{sexp:int * int}\n" (xn, yn); *)
          Queue.enqueue queue (xn, yn))
        ns
    done
  in
  let queue = Queue.singleton s in
  let x, y = s in
  was.(y).(x) <- true;
  impl queue;
  dists

let mark_path map dists s =
  let was = map_array map false in
  let path = map_array map false in
  let rec impl (x, y) =
    if was.(y).(x) then if Cell.is_start map.(y).(x) then true else false
    else (
      (* Stdio.printf !"Enter: %{sexp:int*int}\n" (x, y); *)
      was.(y).(x) <- true;
      let ns = List.filter_map ~f:(check_connected map (x, y)) all_directions in
      let nns = List.filter ~f:impl ns in
      match nns with
      | [] -> false
      | ds ->
          path.(y).(x) <- true;
          true)
  in
  let _ = impl s in
  path

type cloc = A | U | D | N [@@deriving ord]

let cloc_eq lr ll = compare_cloc lr ll = 0

let detect_start_pipe map path (xs, ys) =
  let get_dir_for d =
    let c = check_connected map (xs, ys) d in
    match c with
    | None -> None
    | Some (x, y) -> if path.(y).(x) then Some d else None
  in
  let dirs = List.filter_map ~f:get_dir_for all_directions in
  match dirs with [ f; s ] -> Cell.Pipe (f, s) | _ -> failwith "unreachable"

let get_pipe_type c =
  if Cell.has_diretion T c && Cell.has_diretion B c then A
  else if Cell.has_diretion T c then U
  else if Cell.has_diretion B c then D
  else N

let count_inner map path =
  let res = ref 0 in
  let loc = ref N in
  let nmap = map_array map 0 in
  let h = Array.length map in
  let w = Array.length map.(0) in
  let get_pipe c p =
    if Cell.is_pipe c then c else detect_start_pipe map path p
  in
  for y = 0 to h - 1 do
    loc := N;
    for x = 0 to w - 1 do
      let c = map.(y).(x) in
      if
        Cell.is_ground map.(y).(x)
        || (Cell.is_pipe map.(y).(x) && not path.(y).(x))
      then
        if cloc_eq !loc A then (
          res := !res + 1;
          nmap.(y).(x) <- 1)
        else nmap.(y).(x) <- 2
      else if Cell.is_start c || Cell.is_pipe c then (
        nmap.(y).(x) <- 4;
        if path.(y).(x) then (
          nmap.(y).(x) <- 3;
          let t = get_pipe_type @@ get_pipe c (x, y) in
          match !loc with
          | A -> (
              match t with
              | A -> loc := N
              | U -> loc := D
              | D -> loc := U
              | _ -> ())
          | U -> (
              match t with
              | A -> failwith "unreachable"
              | U -> loc := N
              | D -> loc := A
              | _ -> ())
          | D -> (
              match t with
              | A -> failwith "unreachable"
              | U -> loc := A
              | D -> loc := N
              | _ -> ())
          | N -> (
              match t with
              | A -> loc := A
              | U -> loc := U
              | D -> loc := D
              | _ -> ())))
    done
  done;
  print_map
    (fun s ->
      match s with
      | 0 -> "."
      | 1 -> "I"
      | 2 -> "O"
      | 3 -> "x"
      | 4 -> "_"
      | _ -> failwith "unreachable")
    nmap;
  !res

let () =
  let map = read_input in
  let start = Option.value_exn @@ find_start map in
  (* Stdio.printf !"%{sexp:Cell.t}\n" map.(0).(2); *)
  Stdio.printf !"%{sexp:(int * int) option}\n" @@ check_connected map (1, 2) T;
  Stdio.printf !"Start: %{sexp:int * int}\n" start;
  let dists = bfs map start in
  print_map Int.to_string dists;
  let path = mark_path map dists start in
  print_map (fun b -> if b then "x" else ".") path;
  let res = count_inner map path in
  Stdio.print_endline @@ Int.to_string res;
  (* Stdio.printf !"%d" res; *)
  Stdlib.exit 0
