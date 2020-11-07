str_single_thread = '''func {}() {{
  dieseldist.InitQueues(Num_threads, "amqp://guest:guest@localhost:5672/")
  dieseldist.WaitForWorkers(Num_threads)
  var DynMap [{}]dieseldist.ProbInterval;
  var my_chan_index int;
  _ = my_chan_index;
  _ = DynMap;
  {}

  dieseldist.CleanupMain()
  fmt.Println("Ending thread : ", {});
}}'''

str_single_thread_acc = '''func {}() {{
  dieseldistacc.InitQueues(Num_threads, "amqp://guest:guest@localhost:5672/")
  dieseldistacc.WaitForWorkers(Num_threads)
  var DynMap [{}]float64;
  var my_chan_index int;
  _ = my_chan_index;
  _ = DynMap;
  {}

  dieseldistacc.CleanupMain()
  fmt.Println("Ending thread : ", {});
}}'''

str_single_thread_rel = '''func {}() {{
  dieseldistrel.InitQueues(Num_threads, "amqp://guest:guest@localhost:5672/")
  dieseldistrel.WaitForWorkers(Num_threads)
  var DynMap [{}]float32;
  var my_chan_index int;
  _ = my_chan_index;
  _ = DynMap;
  {}

  dieseldistrel.CleanupMain()
  fmt.Println("Ending thread : ", {});
}}'''

str_member_thread = '''func {}(tid int) {{
  dieseldist.InitQueues(Num_threads, "amqp://guest:guest@localhost:5672/")
  dieseldist.PingMain(tid)
  var DynMap [{}]dieseldist.ProbInterval;
  var my_chan_index int;
  _ = my_chan_index;
  _ = DynMap;
  {}
  fmt.Println("Ending thread : ", {});
}}'''

str_member_thread_acc = '''func {}(tid int) {{ 
  dieseldistacc.InitQueues(Num_threads, "amqp://guest:guest@localhost:5672/")
  dieseldistacc.PingMain(tid)
  var DynMap [{}]float64;
  var my_chan_index int;
  _ = my_chan_index;
  _ = DynMap;
  {}
  fmt.Println("Ending thread : ", {});
}}'''

str_member_thread_rel = '''func {}(tid int) {{
  dieseldistrel.InitQueues(Num_threads, "amqp://guest:guest@localhost:5672/")
  dieseldistrel.PingMain(tid)
  var DynMap [{}]float32;
  var my_chan_index int;
  _ = my_chan_index;
  _ = DynMap;
  {}
  fmt.Println("Ending thread : ", {});
}}'''
