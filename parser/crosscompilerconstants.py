distributed_connection_point = "amqp://guest:guest@localhost:5672/"

single_process_thread = {
    "dieseldist":'''func {}() {{
  dieseldist.InitQueues(Num_threads, "amqp://guest:guest@localhost:5672/")
  dieseldist.WaitForWorkers(Num_threads)
  var DynMap [{}]dieseldist.ProbInterval;
  var my_chan_index int;
  _ = my_chan_index;
  _ = DynMap;
  {}

  // dieseldist.CleanupMain()
  fmt.Println("Ending thread : ", {});
}}''',
    "dieseldistrel":'''func {}() {{
  dieseldistrel.InitQueues(Num_threads, "amqp://guest:guest@localhost:5672/")
  dieseldistrel.WaitForWorkers(Num_threads)
  var DynMap [{}]float32;
  var my_chan_index int;
  _ = my_chan_index;
  _ = DynMap;
  {}

  // dieseldistrel.CleanupMain()
  fmt.Println("Ending thread : ", {});
}}''',    
    "dieseldistacc": '''func {}() {{
  dieseldistacc.InitQueues(Num_threads, "amqp://guest:guest@localhost:5672/")
  dieseldistacc.WaitForWorkers(Num_threads)
  var DynMap [{}]float64;
  var my_chan_index int;
  _ = my_chan_index;
  _ = DynMap;
  {}

  // dieseldistacc.CleanupMain()
  fmt.Println("Ending thread : ", {});
}}'''
}

multiple_process_thread = {
    "dieseldist": '''func {}(tid int) {{
  dieseldist.InitQueues(Num_threads, "amqp://guest:guest@localhost:5672/")
  dieseldist.PingMain(tid)
  var DynMap [{}]dieseldist.ProbInterval;
  var my_chan_index int;
  _ = my_chan_index;
  _ = DynMap;
  {}
  fmt.Println("Ending thread : ", {});
}}''',
    "dieseldistrel":'''func {}(tid int) {{
  dieseldistrel.InitQueues(Num_threads, "amqp://guest:guest@localhost:5672/")
  dieseldistrel.PingMain(tid)
  var DynMap [{}]float32;
  var my_chan_index int;
  _ = my_chan_index;
  _ = DynMap;
  {}
  fmt.Println("Ending thread : ", {});
}}''',
    "dieseldistacc":'''func {}(tid int) {{ 
  dieseldistacc.InitQueues(Num_threads, "amqp://guest:guest@localhost:5672/")
  dieseldistacc.PingMain(tid)
  var DynMap [{}]float64;
  var my_chan_index int;
  _ = my_chan_index;
  _ = DynMap;
  {}
  fmt.Println("Ending thread : ", {});
}}'''
}

access_reliability = {
    "dieseldist":"DynMap[{}].Reliability",
    "dieseldistrel":"DynMap[{}]",
}

prob_assignment_str_const = {
    "dieseldist":"DynMap[{}] = dieseldist.ProbInterval{{{}, 0}};\n",
    "dieseldistrel":"DynMap[{}] = {};\n",
    "dieseldistacc":"DynMap[{}] = 1000000.0;\n"
}

prob_assignment_str_single = {
    "dieseldist":"DynMap[{}].Reliability = DynMap[{}].Reliability * {};\n",
    "dieseldistrel":"DynMap[{}] = DynMap[{}] * {};\n",
    "dieseldistacc":"DynMap[{}] = 1000000.0;\n"
}

dyn_pchoice_str = {
    "dieseldist":"DynMap[{}].Reliability = LIBRARYNAME.Max(0.0, {} - float32({})) * {};",
    "dieseldistrel":"DynMap[{}] = LIBRARYNAME.Max(0.0, {} - float32({})) * {};"
}

dyn_pchoice_accuracy = {
    "dieseldist":"DynMap[{}].Delta = 1000000.0;\n",
    "dieseldistrel":"DynMap[{}] = 1000000.0;\n"
}

convert_cond_str_list = {
    ("float32", "int"): "{} = LIBRARYNAME.DynCondFloat32GeqInt({},{},DynMap[:],{},{},{},{},{},{},{});\n",
    ("float64", "int"): "{} = LIBRARYNAME.DynCondFloat64GeqInt({}, {}, DynMap[:],{},{},{},{},{},{},{});\n",
    ("float32", "float32"): "{} = LIBRARYNAME.DynCondFloat32GeqFloat32({},{},DynMap[:],{},{},{},{},{},{},{});\n",
    ("float64", "float64"): "{} = LIBRARYNAME.DynCondFloat64GeqFloat64({}, {}, DynMap[:], {}, {}, {}, {}, {}, {}, {});\n",
}

condassignment_dyn_str = {
    "dieseldist": '''if temp_bool_{0} != 0 {{DynMap[{1}].Reliability  = DynMap[{2}].Reliability * DynMap[{3}].Reliability}} else {{ DynMap[{1}].Reliability = DynMap[{2}].Reliability * DynMap[{4}].Reliability}};\n''',
    "dieseldistrel": '''if temp_bool_{0} != 0 {{DynMap[{1}]  = DynMap[{2}] * DynMap[{3}]}} else {{ DynMap[{1}] = DynMap[{2}] * DynMap[{4}]}};\n'''
}

dyn_expression_str_single = {
    "dieseldist": "DynMap[{}].Reliability = DynMap[{}].Reliability;\n",
    "dieseldistrel": "DynMap[{}] = DynMap[{}];\n",
    "dieseldistacc": "DynMap[{}] = DynMap[{}];\n"
}

dyn_expression_str_precise = {
    "dieseldist": '''DynMap[{}] = dieseldist.ProbInterval{{1, 0}};\n''',
    "dieseldistrel": '''DynMap[{}] = 1.0;\n''',
    "dieseldistacc": '''DynMap[{}] = 0.0;\n'''
}

dyn_assign_str = {
    "dieseldist": "DynMap[{}].Reliability = {} - {}.0;\n",
    "dieseldistrel": "DynMap[{}] = {} - {}.0;\n"
}

dyn_accuracy_update = {
    "dieseldist": "DynMap[{}].Delta = DynMap[{}].Delta;\n",
    "dieseldistacc": "DynMap[{}] = DynMap[{}];\n"
}

dyn_accuracy_update_double = {
    "dieseldist": "DynMap[{}].Delta = DynMap[{}].Delta + DynMap[{}].Delta;\n",
    "dieseldistacc": "DynMap[{}] = DynMap[{}] + DynMap[{}];\n"
}

dyn_accuracy_mult_single_str = {
    "dieseldist": "DynMap[{0}].Delta = math.Abs(float64({1})) * DynMap[{2}].Delta;\n",
    "dieseldistacc": "DynMap[{0}] = math.Abs(float64({1})) * DynMap[{2}];\n"
}

dyn_accuracy_mult_double_str = {
    "dieseldist": "DynMap[{0}].Delta = math.Abs(float64({1})) * DynMap[{2}].Delta + math.Abs(float64({3})) * DynMap[{4}].Delta + DynMap[{2}].Delta*DynMap[{4}].Delta;\n",
    "dieseldistacc": "DynMap[{0}] = math.Abs(float64({1})) * DynMap[{2}] + math.Abs(float64({3})) * DynMap[{4}] + DynMap[{2}]*DynMap[{4}];\n"
}

dyn_accuracy_div_single_str_0 = {
    "dieseldist": "DynMap[{0}].Delta =  DynMap[{2}].Delta / math.Abs(float64({1}));\n",
    "dieseldistacc": "DynMap[{0}] =  DynMap[{2}] / math.Abs(float64({1}));\n"
}

dyn_accuracy_div_single_str_1 = {
    "dieseldist": "DynMap[{0}].Delta =  DynMap[{2}].Delta * math.Abs(float64({1}));\n",
    "dieseldistacc": "DynMap[{0}] =  DynMap[{2}] * math.Abs(float64({1}));\n"
}

dyn_accuracy_div_double_str = {
    "dieseldist": "DynMap[{0}].Delta = math.Abs(float64({1})) * DynMap[{2}].Delta + math.Abs(float64({3})) * DynMap[{4}].Delta / (math.Abs(float64({3})) * (math.Abs(float64({1}))-DynMap[{4}].Delta));\n",
    "dieseldistacc": "DynMap[{0}] = math.Abs(float64({1})) * DynMap[{2}] + math.Abs(float64({3})) * DynMap[{4}] / (math.Abs(float64({3})) * (math.Abs(float64({1}))-DynMap[{4}]));\n"
}

dyn_array_store_precise = {
    "dieseldist": "DynMap[{} + _temp_index_{}] = dieseldist.ProbInterval{{1, 0}};\n",
    "dieseldistrel": "DynMap[{} + _temp_index_{}] = 1.0;\n",
    "dieseldistacc": "DynMap[{} + _temp_index_{}] = 0.0;\n",
}

dyn_array_store_single_dyn = "DynMap[{0} + _temp_index_{2}] = DynMap[{1}];\n"

dyn_cast_dyn_str = {
    "dieseldist": "DynMap[{0}].Reliability = DynMap[{3}].Reliability;\n DynMap[{0}].Delta = dieseldist.GetCastingError64to32({1}, {2});\n",
    "dieseldistrel": "DynMap[{0}] = DynMap[{3}];\n",
    "dieseldistacc": "DynMap[{0}] = dieseldistacc.GetCastingError64to32({1}, {2});\n",
}

dyn_cast_precise_str = {
    "dieseldist": "DynMap[{0}].Reliability=1;\nDynMap[{0}].Delta=dieseldist.GetCastingError64to32({1}, {2});\n",
    "dieseldistrel": "DynMap[{0}] = 1.0;\n",
    "dieseldistacc": "DynMap[{0}] = dieseldistacc.GetCastingError64to32({1}, {2});\n",
}

dyn_track_str = {
    "dieseldist": "DynMap[{0}] = dieseldist.ProbInterval{{{1}, {2}}};\n",
    "dieseldistrel": "DynMap[{0}] = {1};\n",
    "dieseldistacc": "DynMap[{0}] = {2};\n",
}

dyn_init_str = {
    "dieseldist": "DynMap[{}] = dieseldist.ProbInterval{{1, 0}};\n",
    "dieseldistrel": "DynMap[{}] = 1.0;\n",
    "dieseldistacc": "DynMap[{}] = 0.0;\n",
}

str_probchoiceIntFlag = "{} = LIBRARYNAME.RandchoiceFlag(float32({}), {}, {}, &__flag_{});\n"
str_probchoiceInt = "{} = LIBRARYNAME.Randchoice(float32({}), {}, {});\n"

dyn_rec_str = '''my_chan_index = {0} * LIBRARYNAME.Numprocesses + {1};
__temp_rec_val_{3} := LIBRARYNAME.GetDynValue(my_chan_index);
DynMap[{2}] = __temp_rec_val_{3};
'''

ch_str = '''
fmt.Println("----------------------------");\n
fmt.Println("Spec checkarray({3}, {1}): ", LIBRARYNAME.CheckArray({0}, {1}, {2}, DynMap[:]));\n
fmt.Println("----------------------------");\n
'''

ch_str_result = "{0} = LIBRARYNAME.Check(DynMap[{1}], {2}, {3});\n"

cond_send_str = {
    ("int", 0): "LIBRARYNAME.Condsend({}, {}, {}, {});\n",
    ("int32", 0): "LIBRARYNAME.CondsendInt32({}, {}, {}, {});\n",
    ("int64", 0): "LIBRARYNAME.CondsendInt64({}, {}, {}, {});\n",
    ("float32", 0): "LIBRARYNAME.CondsendFloat32({}, {}, {}, {});\n",
    ("float64", 0): "LIBRARYNAME.CondsendFloat64({}, {}, {}, {});\n",
    ("int", 1): "LIBRARYNAME.CondsendIntArray({}, {}[:], {}, {});\n",
    ("int", 1): "LIBRARYNAME.CondsendIntArray({}, {}[:], {}, {});\n",
}

send_str = {
    ("int", 0): "LIBRARYNAME.SendInt({}, {}, {});\n",
    ("int32", 0): "LIBRARYNAME.SendInt32({}, {}, {});\n",
    ("int64", 0): "LIBRARYNAME.SendInt64({}, {}, {});\n",
    ("float32", 0): "LIBRARYNAME.SendFloat32({}, {}, {});\n",
    ("float64", 0): "LIBRARYNAME.SendFloat64({}, {}, {});\n",
    ("int", 1): "LIBRARYNAME.SendIntArray({}[:], {}, {});\n",
    ("int32", 1): "LIBRARYNAME.SendInt32Array({}[:], {}, {});\n",
    ("int64", 1): "LIBRARYNAME.SendInt64Array({}[:], {}, {});\n",
    ("float32", 1): "LIBRARYNAME.SendFloat32Array({}[:], {}, {});\n",
    ("float64", 1): "LIBRARYNAME.SendFloat64Array({}[:], {}, {});\n"
}

dyn_send_str = {
    ("int", 0): "LIBRARYNAME.SendDynIntArray({}[:], {}, {}, DynMap[:], {});\n",
    ("float64", 0): "LIBRARYNAME.SendDynFloat64Array({}[:], {}, {}, DynMap[:], {});\n",
    ("float32", 0): "LIBRARYNAME.SendDynFloat32Array({}[:], {}, {}, DynMap[:], {});\n",
    ("int", 1): "LIBRARYNAME.SendDynIntArrayO1({}[:], {}, {}, DynMap[:], {});\n",
    ("float64", 1): "LIBRARYNAME.SendDynFloat64ArrayO1({}[:], {}, {}, DynMap[:], {});\n",
    ("float32", 1): "LIBRARYNAME.SendDynFloat32ArrayO1({}[:], {}, {}, DynMap[:], {});\n",
}

rec_str = {
    ("int", 0): "LIBRARYNAME.ReceiveInt(&{}, {}, {});\n",
    ("int32", 0): "LIBRARYNAME.ReceiveInt32(&{}, {}, {});\n",
    ("int64", 0): "LIBRARYNAME.ReceiveInt64(&{}, {}, {});\n",
    ("float32", 0): "LIBRARYNAME.ReceiveFloat32(&{}, {}, {});\n",
    ("float64", 0): "LIBRARYNAME.ReceiveFloat64(&{}, {}, {});\n",
    ("int", 1): "LIBRARYNAME.ReceiveIntArray({}[:], {}, {});\n",
    ("int32", 1): "LIBRARYNAME.ReceiveInt32Array({}[:], {}, {});\n",
    ("int64", 1): "LIBRARYNAME.ReceiveInt64Array({}[:], {}, {});\n",
    ("float32", 1): "LIBRARYNAME.ReceiveFloat32Array({}[:], {}, {});\n",
    ("float64", 1): "LIBRARYNAME.ReceiveFloat64Array({}[:], {}, {});\n"
}

dyn_rec_dict = {
    ("int", 0): "LIBRARYNAME.ReceiveDynIntArray({}[:], {}, {}, DynMap[:], {});\n",
    ("float64", 0): "LIBRARYNAME.ReceiveDynFloat64Array({}[:], {}, {}, DynMap[:], {});\n",
    ("float32", 0): "LIBRARYNAME.ReceiveDynFloat32Array({}[:], {}, {}, DynMap[:], {});\n",
    ("int", 1): "LIBRARYNAME.ReceiveDynIntArrayO1({}[:], {}, {}, DynMap[:], {});\n",
    ("float64", 1): "LIBRARYNAME.ReceiveDynFloat64ArrayO1({}[:], {}, {}, DynMap[:], {});\n",
    ("float32", 1): "LIBRARYNAME.ReceiveDynFloat32ArrayO1({}[:], {}, {}, DynMap[:], {});\n",
}

noisy_dyn_rec_dict = {
    ("int", 0): "LIBRARYNAME.NoisyReceiveDynIntArray({}[:], {}, {}, DynMap[:], {});\n",
    ("float64", 0): "LIBRARYNAME.NoisyReceiveDynFloat64Array({}[:], {}, {}, DynMap[:], {});\n",
    # ("float32", 0): "LIBRARYNAME.NoisyReceiveDynFloat32Array({}[:], {}, {}, DynMap[:], {});\n",
            ("int", 1): "LIBRARYNAME.NoisyReceiveDynIntArrayO1({}[:], {}, {}, DynMap[:], {});\n",
    ("float64", 1): "LIBRARYNAME.NoisyReceiveDynFloat64ArrayO1({}[:], {}, {}, DynMap[:], {});\n",
    # ("float32", 1): "LIBRARYNAME.NoisyReceiveDynFloat32ArrayO1({}[:], {}, {}, DynMap[:], {});\n",
}
