Q = {1,2,3,4,5,6,7,8,9,10};

0:[
  precise int32 [10][10] edges;
  precise int32 [10] inlinks;
  precise int32 [10] outlinks;
  dynamic float64 [10] pageranks;
  precise int32 q;
  precise int32 nodeInlinks;
  dynamic float64 newPagerank;
  dynamic int32 receivedPagerank;
  dynamic float64 temp;

  for q in Q do {
    send(q, precise int32 [], edges);
    nodeInlinks = inlinks[q-1];
    send(q, precise int32, nodeInlinks);
    send(q, precise int32 [], outlinks);
  };
  for q in Q do {
      dynsend(q, dynamic float64[], pageranks);
  };
  for q in Q do {
      receivedPagerank, newPagerank = dyncondreceive(q, dynamic float64);
      temp = pageranks[q-1];
      temp = newPagerank [receivedPagerank] temp;
      pageranks[q-1] = temp;
  };
]

||

q in Q:[
  precise int32 [10][10] edges;
  precise int32 inlinks;
  precise int32 [10] outlinks;
  dynamic int32 [10] outlinksd;
  dynamic float64 [10] pageranks;
  dynamic float64 newPagerank;
  precise int32 inlink;
  precise int32 neighbor;
  dynamic  int32 toSendOrNot;
  dynamic float64 delta;

  edges = receive(0, precise int32[]);
  inlinks = receive(0, precise int32);
  outlinks = receive(0, precise int32[]);

  outlinksd = track(outlinks, 1.0);

  pageranks = dynreceive(0, dynamic float64[]);
  newPagerank = 0.15;
  inlink = 0;
  while (inlink < inlinks)  {
        neighbor = edges[q-1][inlink];
        newPagerank = newPagerank + 0.85*pageranks[neighbor]/outlinksd[neighbor];
        inlink = inlink+1;
  };
  delta = pageranks[q-1]-newPagerank;
  toSendOrNot = (delta>0.01);
  dyncondsend(toSendOrNot, 0, dynamic float64, newPagerank);
]