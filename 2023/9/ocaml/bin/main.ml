open Base

let read_input =
  let contents = Stdio.In_channel.read_lines "input.txt" in
  let seqs =
    List.map
      ~f:(fun s ->
        List.map ~f:Int.of_string @@ String.split ~on:' ' s)
      contents
  in
  seqs

let get_diffs seq =
  let rec impl res = function
    | [] -> res
    | [ _ ] -> res
    | f :: s :: rest -> impl ((f - s) :: res) (s :: rest)
  in
  List.rev @@ impl [] seq

let pred_next seq =
  let rec impl cur =
    if List.for_all ~f:(( = ) 0) cur then 0 :: cur
    else
      let diffs = get_diffs cur in
      let cont_diffs = impl diffs in
      let cont_cur =
        match (cont_diffs, cur) with
        | d :: _, n :: _ -> (d + n) :: cur
        | _ -> failwith "unreachable"
      in
      (* Stdio.printf *)
      (* !"Diffs: %{sexp:int list}, Cont diffs: %{sexp:int list}, Cur: \ *)
         (*     %{sexp:int list}, Cont Cur: %{sexp:int list}\n" *)
      (*   diffs cont_diffs cur cont_cur; *)
      cont_cur
  in
  impl seq

let () =
  let inp = read_input in
  Stdio.printf !"%{sexp:int list list}\n" inp;
  let cont = List.map ~f:pred_next inp in
  Stdio.printf !"%{sexp:int list list}\n" cont;
  let res = List.fold ~init:0 ~f:(fun acc l -> acc + List.hd_exn l) cont in
  Stdio.print_endline @@ Int.to_string res
