1:[precise int x; approx int y; precise bool p; x=x*10; p=x==10; if p then {x=x/10} else {y=x+y; skip}; y=receive(2, approx int); y=receive(3, approx int)]
