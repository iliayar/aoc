open Batteries

module Field = struct
  type cell = Empty | Num of int | Symbol of char | Gear
  type t = cell array array

  let is_symbol = function Symbol _ -> true | _ -> false
  let is_gear = function Gear -> true | _ -> false

  let cell_of_char = function
    | '0' .. '9' as c -> Num (Char.code c - Char.code '0')
    | '.' -> Empty
    | '*' -> Gear
    | c -> Symbol c

  let of_list : _ -> t =
    Array.map (Array.map cell_of_char % Array.of_list) % Array.of_list

  let shape (f : t) =
    let rows = Array.length f in
    let columns = Array.length f.(0) in
    (rows, columns)

  let get (f : t) (x, y) =
    let rows, cols = shape f in
    if x < 0 || x >= rows then None
    else if y < 0 || y >= cols then None
    else Some f.(x).(y)
end

module CoordsCmpTrait = struct
  type t = int * int

  let compare (xl, yl) (xr, yr) =
    if xl < xr then -1
    else if xl > xr then 1
    else if yl < yr then -1
    else if yl > yr then 1
    else 0
end

module CoordsMap = Map.Make (CoordsCmpTrait)
module CoordsSet = Set.Make (CoordsCmpTrait)

type num_segment = { n : int; coords : (int * int) list }

let find_nums s =
  let rec impl s cur_n idx idxs =
    match s with
    | ('0' .. '9' as c) :: ss ->
        impl ss
          ((cur_n * 10) + (Char.code c - Char.code '0'))
          (idx + 1) (idx :: idxs)
    | s ->
        let rest_items =
          match s with _ :: ss -> impl ss 0 (idx + 1) [] | [] -> []
        in
        if List.is_empty idxs then rest_items else (cur_n, idxs) :: rest_items
  in
  impl s 0 0 []

let read_input =
  let content = File.lines_of "input.txt" in
  let content = List.map String.to_list @@ List.of_enum content in
  let nums =
    List.flatten
    @@ List.mapi
         (fun i l ->
           List.map (fun (n, idxs) ->
               let coords = List.map (fun j -> (i, j)) idxs in
               { n; coords })
           @@ find_nums l)
         content
  in
  (Field.of_list content, nums)

let get_coords_around (x, y) =
  [
    (x - 1, y - 1);
    (x, y - 1);
    (x + 1, y - 1);
    (x - 1, y);
    (x + 1, y);
    (x - 1, y + 1);
    (x, y + 1);
    (x + 1, y + 1);
  ]

let find_gears_around f =
  List.filter_map (fun c ->
      Option.bind (Field.get f c) (fun cell ->
          if Field.is_gear cell then Some c else None))
  % get_coords_around

let find_gears_around_num f { coords; _ } =
  let found = List.map (find_gears_around f) coords in
  List.fold
    (fun acc cs -> CoordsSet.union acc @@ CoordsSet.of_list cs)
    CoordsSet.empty found

let update_gears_for_num f m num =
  let coords_set = find_gears_around_num f num in
  CoordsSet.fold
    CoordsMap.(
      fun c m ->
        let cl = find_opt c m |? [] in
        add c (num.n :: cl) m)
    coords_set m

let find_gears f = List.fold (update_gears_for_num f) CoordsMap.empty

let solve =
  let f, nums = read_input in
  let gears = find_gears f nums in
  let gears_power =
    CoordsMap.filter_map
      (fun _ ns -> match ns with [ n1; n2 ] -> Some (n1 * n2) | _ -> None)
      gears
  in
  let r = Enum.sum @@ CoordsMap.values gears_power in
  print_endline @@ Int.to_string r

let () = solve
