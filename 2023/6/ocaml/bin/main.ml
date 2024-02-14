open Base

let read_input =
  let lines = Stdio.In_channel.read_lines "input.txt" in
  let time_str, distance_str =
    match lines with
    | [ time; distance ] ->
        ( String.chop_prefix_exn ~prefix:"Time:" time,
          String.chop_prefix_exn ~prefix:"Distance:" distance )
    | _ -> failwith "unreachable"
  in
  let cnv s =
    Int.of_string @@ String.concat
    @@ List.filter ~f:(Fn.compose not String.is_empty)
    @@ String.split ~on:' ' s
  in
  let data = [((cnv time_str), (cnv distance_str))] in
  data

let solve' (t, d) =
  let t = Float.of_int t in
  let d = Float.of_int d in
  let discr = (t **. 2.) -. (4. *. d) in
  let sd = Float.sqrt discr in
  let l = (t -. sd) /. 2. in
  let r = (t +. sd) /. 2. in
  let li = Int.of_float @@ Float.round_down l in
  let ri = Int.of_float @@ Float.round_up r in
  ri - li - 1

let solve races = List.fold ~init:1 ~f:(fun acc r -> acc * solve' r) races

let () =
  let input = read_input in
  Stdio.print_endline @@ Int.to_string @@ solve input
