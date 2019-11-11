  var sign int
  var OutputX DynFloat32
  var xInput DynFloat32
  var xNPrimeofX DynFloat32
  var expValues DynFloat32
  var xK2 DynFloat32
  var xK2_2 DynFloat32
  var xK2_3 DynFloat32
  var xK2_4 DynFloat32
  var xK2_5 DynFloat32
  var xLocal DynFloat32
  var xLocal_1 DynFloat32
  var xLocal_2 DynFloat32
  var xLocal_3 DynFloat32
  if InputX.Num < 0.0 {
    InputX = MulDynFloat32(InputX,MakeDynFloat32(-1))
    sign = 1
  } else {
    sign = 0
  }
  xInput = InputX
  expValues = ExpDynFloat32(MulDynFloat32(MakeDynFloat32(-0.5),MulDynFloat32(InputX,InputX)))
  xNPrimeofX = expValues
  xNPrimeofX = MulDynFloat32(xNPrimeofX,MakeDynFloat32(0.39894228040143270286))
  xK2 = MulDynFloat32(MakeDynFloat32(0.2316419),xInput)
  xK2 = AddDynFloat32(MakeDynFloat32(1.0),xK2)
  xK2 = DivDynFloat32(MakeDynFloat32(1.0),xK2)
  xK2_2 = MulDynFloat32(xK2,xK2)
  xK2_3 = MulDynFloat32(xK2_2,xK2)
  xK2_4 = MulDynFloat32(xK2_3,xK2)
  xK2_5 = MulDynFloat32(xK2_4,xK2)
  xLocal_1 = MulDynFloat32(xK2,MakeDynFloat32(0.319381530))
  xLocal_2 = MulDynFloat32(xK2_2,MakeDynFloat32(-0.356563782))
  xLocal_3 = MulDynFloat32(xK2_3,MakeDynFloat32(1.781477937))
  xLocal_2 = AddDynFloat32(xLocal_2,xLocal_3)
  xLocal_3 = MulDynFloat32(xK2_4,MakeDynFloat32(-1.821255978))
  xLocal_2 = AddDynFloat32(xLocal_2,xLocal_3)
  xLocal_3 = MulDynFloat32(xK2_5,MakeDynFloat32(1.330274429))
  xLocal_2 = AddDynFloat32(xLocal_2,xLocal_3)
  xLocal_1 = AddDynFloat32(xLocal_2,xLocal_1)
  xLocal   = MulDynFloat32(xLocal_1,xNPrimeofX)
  xLocal   = SubDynFloat32(MakeDynFloat32(1.0),xLocal)
  OutputX  = xLocal
  if sign==1 {
    OutputX = SubDynFloat32(MakeDynFloat32(1.0),OutputX)
  }
  return OutputX

  var OptionPrice DynFloat32
  var xStockPrice DynFloat32
  var xStrikePrice DynFloat32
  var xRiskFreeRate DynFloat32
  var xVolatility DynFloat32
  var xTime DynFloat32
  var xSqrtTime DynFloat32
  var logValues DynFloat32
  var xLogTerm DynFloat32
  var xD1 DynFloat32
  var xD2 DynFloat32
  var xPowerTerm DynFloat32
  var xDen DynFloat32
  var d1 DynFloat32
  var d2 DynFloat32
  var FutureValueX DynFloat32
  var NofXd1 DynFloat32
  var NofXd2 DynFloat32
  var NegNofXd1 DynFloat32
  var NegNofXd2 DynFloat32
  xStockPrice = sptprice
  xStrikePrice = strike
  xRiskFreeRate = rate
  xVolatility = volatility
  xTime = time
  xSqrtTime = SqrtDynFloat32(xTime)
  logValues = LogDynFloat32(DivDynFloat32(xStockPrice,xStrikePrice))
  xLogTerm = logValues
  xPowerTerm = MulDynFloat32(xVolatility,xVolatility)
  xPowerTerm = MulDynFloat32(xPowerTerm,MakeDynFloat32(0.5))
  xD1 = AddDynFloat32(xRiskFreeRate,xPowerTerm)
  xD1 = MulDynFloat32(xD1,xTime)
  xD1 = AddDynFloat32(xD1,xLogTerm)
  xDen = MulDynFloat32(xVolatility,xSqrtTime)
  xD1 = DivDynFloat32(xD1,xDen)
  xD2 = SubDynFloat32(xD1,xDen)
  d1 = xD1
  d2 = xD2
  NofXd1 = CNDF( d1 )
  NofXd2 = CNDF( d2 )
  FutureValueX = MulDynFloat32(xStrikePrice,ExpDynFloat32(MulDynFloat32(MulDynFloat32(rate,time),MakeDynFloat32(-1))))
  if (otype.Num == 0) {
    OptionPrice = SubDynFloat32(MulDynFloat32(xStockPrice,NofXd1),MulDynFloat32(FutureValueX,NofXd2))
  } else {
    NegNofXd1 = SubDynFloat32(MakeDynFloat32(1.0),NofXd1)
    NegNofXd2 = SubDynFloat32(MakeDynFloat32(1.0),NofXd2)
    OptionPrice = SubDynFloat32(MulDynFloat32(FutureValueX,NegNofXd2),MulDynFloat32(xStockPrice,NegNofXd1))
  }
  return OptionPrice

  stockstrings := <- chin
  var mystocks []DynFloat32
  for i := 0; i < len(stockstrings); i++ {
    c := stockstrings[i]
    optionPrice := BlkSchlsEqEuroNoDiv(c[0], c[1], c[2], c[4], c[5], c[6], c[8])
    mystocks = append(mystocks, optionPrice)
  }
  for i := 0; i < len(mystocks); i++ {
    chout <- DynFloat32(mystocks[i])
  }

  var data_array []([] DynFloat32)
  var results []DynFloat32
  for i := 0; i < num_threads; i++ {
    cinput[i] <- data_array[workperthread*i:workperthread*(i+1)]
  }
  for i := 0; i < num_threads; i++ {
    for j:=0; j < workperthread; j++ {
      result := <- coutput[i]
      results = append(results, DynFloat32(result))
    }
  }
