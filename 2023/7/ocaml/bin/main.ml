open Base

module Card = struct
  type t = J | N of int | Q | K | A [@@deriving ord, sexp]

  let of_char c =
    match c with
    | 'A' -> A
    | 'K' -> K
    | 'Q' -> Q
    | 'J' -> J
    | 'T' -> N 10
    | '2' .. '9' -> N (Char.to_int c - Char.to_int '0')
    | _ -> failwith "unreachable"
end

module Hand = struct
  type t = Card.t list [@@deriving ord, sexp]

  let of_string s : t =
    if String.length s <> 5 then failwith "unreachable";
    List.map ~f:Card.of_char @@ String.to_list s

  type counts_t = (Card.t * int) list * int

  let counts (hand : t) : counts_t =
    let rec put' assoc card =
      match assoc with
      | ((c, n) as e) :: rest ->
          if Card.compare card c = 0 then (card, n + 1) :: rest
          else e :: put' rest card
      | [] -> [ (card, 1) ]
    in
    let put (assoc, jcnt) card =
      if Card.compare Card.J card = 0 then (assoc, jcnt + 1)
      else (put' assoc card, jcnt)
    in
    let counts, jcnt =
      List.fold ~init:([], 0) ~f:(fun acc card -> put acc card) hand
    in
    let cmp (lc, ln) (rc, rn) =
      if ln = rn then Card.compare rc lc else Int.compare rn ln
    in
    (List.sort ~compare:cmp counts, jcnt)
end

module Combination = struct
  type t = None | Pair | TwoPairs | Three | FullHouse | Four | Five
  [@@deriving ord, sexp]

  let of_hand_counts ((counts, jcnt) : Hand.counts_t) =
    match counts with
    | (_, 5) :: _ -> Five
    | (_, 4) :: _ -> (
        match jcnt with 0 -> Four | 1 -> Five | _ -> failwith "unreachable")
    | (_, 3) :: (_, 2) :: _ -> FullHouse
    | (_, 3) :: _ -> (
        match jcnt with
        | 0 -> Three
        | 1 -> Four
        | 2 -> Five
        | _ -> failwith "unreachable")
    | (_, 2) :: (_, 2) :: _ -> (
        match jcnt with
        | 0 -> TwoPairs
        | 1 -> FullHouse
        | _ -> failwith "unreachable")
    | (_, 2) :: _ -> (
        match jcnt with
        | 0 -> Pair
        | 1 -> Three
        | 2 -> Four
        | 3 -> Five
        | _ -> failwith "unreachable")
    | _ -> (
        match jcnt with
        | 0 -> None
        | 1 -> Pair
        | 2 -> Three
        | 3 -> Four
        | 4 -> Five
        | 5 -> Five
        | _ -> failwith "unreachable")

  let of_hand (hand : Hand.t) = of_hand_counts @@ Hand.counts hand
end

let compare_hands (lcomb, lhand) (rcomb, rhand) =
  let cc = Combination.compare lcomb rcomb in
  if cc = 0 then Hand.compare lhand rhand else cc

let read_input =
  let content = Stdio.In_channel.read_lines "input.txt" in
  let parse_line l =
    let hand, bid =
      match String.split ~on:' ' l with
      | [ hand; bid ] -> (hand, bid)
      | _ -> failwith "unreachable"
    in
    (Hand.of_string hand, Int.of_string bid)
  in
  List.map ~f:parse_line content

let () =
  let input = read_input in
  Stdio.printf !"%{sexp:(Hand.t * int) list}\n" input;
  (* let counts = List.map ~f:(fun (hand, _) -> Hand.counts hand) input in *)
  (* Stdio.printf !"Counts:\n%{sexp:(Card.t * int) list list}\n" counts; *)
  (* let combinations = List.map ~f:Combination.of_hand_counts counts in *)
  (* Stdio.printf !"Combinations:\n%{sexp:Combination.t list}\n" combinations; *)
  let data =
    List.map ~f:(fun (hand, bid) -> (Combination.of_hand hand, hand, bid)) input
  in
  let sorted =
    List.sort
      ~compare:(fun (lcomb, lhand, _) (rcomb, rhand, _) ->
        compare_hands (lcomb, lhand) (rcomb, rhand))
      data
  in
  Stdio.printf !"Result:\n%{sexp:(Combination.t * Hand.t * int) list}\n" sorted;
  let result =
    List.foldi
      ~f:(fun i acc (_, _, bid) -> (bid * (i + 1)) + acc)
      ~init:0 sorted
  in
  Stdio.print_endline @@ Int.to_string result
