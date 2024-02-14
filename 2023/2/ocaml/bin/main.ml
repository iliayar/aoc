open Batteries

type cubes_color = Green | Red | Blue
    (* [@@deriving show] *)

type cubes_set = (cubes_color * int) list
    (* [@@deriving show] *)

type game = 
    { num : int;
      sets : cubes_set list;
    }
    (* [@@deriving show] *)

let read_games = 
    let parse_color = function
        | "green" -> Green
        | "red" -> Red
        | "blue" -> Blue
        | color -> failwith ("unknown color: " ^ color)
    in

    let parse_pair str =
        let (n_str, c_str) = String.split ~by:" " str in
        (parse_color c_str, Int.of_string n_str)
    in

    let parse_set str =
        let pairs = String.split_on_string ~by:", " str in
        List.map parse_pair pairs
    in

    let parse_sets str = 
        let sets = String.split_on_string ~by:"; " str in
        List.map parse_set sets 
    in

    let parse_line line =
        let parts = 
            let r = Str.regexp {|Game \([0-9]+\): \(.*\)|} in
            if Str.string_match r line 0 then
                let n_str = Str.matched_group 1 line in
                let sets_str = Str.matched_group 2 line in
                { num = Int.of_string n_str;
                  sets = parse_sets sets_str;
                }
            else failwith "unreachable"
        in
        parts
    in

    let lines = File.lines_of "input.txt" in
    Enum.map parse_line lines

let get_max_num ~color sets =
    List.fold (fun acc set -> max acc (List.assoc_opt color set |? 0)) 0 sets

let get_min_amount sets = 
    let colors = [Green; Red; Blue] in
    List.map (fun color -> (color, get_max_num ~color sets)) colors

let get_power sets =
    let min_amount = get_min_amount sets in
    List.fold (fun acc (_, n) -> acc * n) 1 min_amount

let solve games = 
    let powers = Enum.map (fun game -> get_power game.sets) games in
    Enum.fold (+) 0 powers
                
let () = 
    let games = read_games in
    let res = solve @@ games in
    res |> Int.to_string |> print_endline
