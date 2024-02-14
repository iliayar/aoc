open Base

type card = { n : int; my_list : int list; win_list : int list }
[@@deriving sexp]

let split ~on s =
  List.filter ~f:(fun s -> not @@ String.is_empty s) @@ String.split ~on s

let read_input =
  let lines = Stdio.In_channel.read_lines "input.txt" in
  let parse_line line =
    match String.split ~on:':' line with
    | [ code_n_s; lists_s ] ->
        let n =
          match split ~on:' ' code_n_s with
          | [ _; n ] -> Int.of_string n
          | _ -> failwith "unreachable"
        in
        let act, win =
          match String.split ~on:'|' lists_s with
          | [ act_s; win_s ] ->
              let conv_list s = List.map ~f:Int.of_string @@ split ~on:' ' s in
              let act = conv_list act_s in
              let win = conv_list win_s in
              (act, win)
          | _ -> failwith "unreachable"
        in
        { n; my_list = act; win_list = win }
    | _ -> failwith "unreachable"
  in
  List.map ~f:parse_line lines

let find_wins card =
  List.filter
    ~f:(fun n -> Option.is_some @@ List.find ~f:(Int.equal n) card.win_list)
    card.my_list

let get_score card =
  let wins = find_wins card in
  let wins_cnt = List.length wins in
  if Int.equal wins_cnt 0 then 0 else Int.pow 2 @@ (wins_cnt - 1)

let get_scores cards = List.map ~f:get_score cards

let solve cards =
  let counts = Array.create ~len:(1 + List.length cards) 1 in
  counts.(0) <- 0;
  let inc_next c n =
    let factor = counts.(c) in
    for i = c + 1 to c + n do
      counts.(i) <- counts.(i) + (1 * factor)
    done
  in
  List.iter
    ~f:(fun card ->
      let wins = find_wins card in
      inc_next card.n @@ List.length wins)
    cards;
  counts

let () =
  let cards = read_input in
  let counts = solve cards in
  let res = Array.fold counts ~init:0 ~f:Int.( + ) in
  Stdio.print_endline @@ Int.to_string res
