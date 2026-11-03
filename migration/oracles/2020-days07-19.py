"""Independent answer calculation: pass the original repository checkout as argv[1]."""
from pathlib import Path
import sys
import json,hashlib,re,itertools,math,functools,collections,time
ROOT=Path(sys.argv[1])/'2020'/'inputs';OUT={}
def save(d,a,b,extra=None):
 p=ROOT/f'day{d:02}.txt';OUT[str(d)]={'part1':str(a),'part2':str(b),'input_path':str(p),'input_sha256':hashlib.sha256(p.read_bytes()).hexdigest() if p.exists() else None}
 if extra:OUT[str(d)].update(extra)
 print(d,a,b,flush=True)
def text(d):return (ROOT/f'day{d:02}.txt').read_text().strip()
# Bag graph: ancestor closure and multiplicity recurrence.
rules={}
for l in text(7).splitlines():
 head,tail=l.split(' bags contain ');rules[head]=[(int(n),c) for n,c in re.findall(r'(\d+) (.+?) bags?',tail)]
parents={'shiny gold'}
while True:
 more={k for k,v in rules.items() if any(c in parents for n,c in v)}-parents
 if not more:break
 parents|=more
@functools.cache
def count(k):return sum(n*(1+count(c)) for n,c in rules[k])
save(7,len(parents)-1,count('shiny gold'))
# VM execution and exhaustive one-instruction repair.
code=[(l.split()[0],int(l.split()[1])) for l in text(8).splitlines()]
def run(flip=-1):
 ip=acc=0;seen=set()
 while 0<=ip<len(code) and ip not in seen:
  seen.add(ip);op,arg=code[ip]
  if ip==flip:op={'jmp':'nop','nop':'jmp'}[op]
  if op=='acc':acc+=arg
  ip+=arg if op=='jmp' else 1
 return ip,acc
p1=run()[1];fixed=[run(i)[1] for i,(o,a) in enumerate(code) if o!='acc' and run(i)[0]==len(code)]
assert len(fixed)==1;save(8,p1,fixed[0])
# Preamble and all contiguous ranges bounded by target (all values are positive).
nums=list(map(int,text(9).splitlines()));bad=next(x for i,x in enumerate(nums[25:],25) if not any(a!=b and a+b==x for a,b in itertools.combinations(nums[i-25:i],2)))
windows=[]
for start in range(len(nums)):
 total=0
 for end in range(start,len(nums)):
  total+=nums[end]
  if total==bad and end>start:windows.append(nums[start:end+1])
  if total>=bad:break
assert len(windows)==1;save(9,bad,min(windows[0])+max(windows[0]))
# Sorted adapter path count.
j=[0]+sorted(map(int,text(10).splitlines()));j+=[j[-1]+3];diff=collections.Counter(b-a for a,b in zip(j,j[1:]));ways={0:1}
for x in j[1:]:ways[x]=sum(ways.get(x-step,0) for step in (1,2,3))
save(10,diff[1]*diff[3],ways[j[-1]])
# Seat sets and simultaneous updates; precompute sightlines independently of state.
grid=text(11).splitlines();h=len(grid);w=len(grid[0]);seats={(x,y) for y,row in enumerate(grid) for x,c in enumerate(row) if c!='.'};deltas=[(x,y) for x in (-1,0,1) for y in (-1,0,1) if x or y]
def stable(visible):
 adj={}
 for x,y in seats:
  ns=[]
  for dx,dy in deltas:
   xx,yy=x+dx,y+dy
   while visible and 0<=xx<w and 0<=yy<h and (xx,yy) not in seats:xx+=dx;yy+=dy
   if (xx,yy) in seats:ns.append((xx,yy))
  adj[x,y]=ns
 occ={(x,y) for x,y in seats if grid[y][x]=='#'}
 while True:
  nxt={p for p in seats if (sum(n in occ for n in adj[p])< (5 if visible else 4) if p in occ else not any(n in occ for n in adj[p]))}
  if nxt==occ:return len(occ)
  occ=nxt
save(11,stable(False),stable(True))
# Navigation: complex numbers, quarter-turn multiplication.
ship=0j;heading=1+0j;boat=0j;waypoint=10+1j
for l in text(12).splitlines():
 op,n=l[0],int(l[1:])
 if op in 'NSEW':
  v={'N':1j,'S':-1j,'E':1,'W':-1}[op]*n;ship+=v;waypoint+=v
 elif op in 'LR':
  rot=(1j if op=='L' else -1j)**(n//90);heading*=rot;waypoint*=rot
 else:ship+=heading*n;boat+=waypoint*n
save(12,int(abs(ship.real)+abs(ship.imag)),int(abs(boat.real)+abs(boat.imag)))
# CRT merge using modular inverse, allowing compatible shared factors.
a,b=text(13).splitlines();earliest=int(a);buses=[(i,int(x)) for i,x in enumerate(b.split(',')) if x!='x'];wait,bus=min(((-earliest)%n,n) for i,n in buses);t=0;mod=1
for i,n in buses:
 g=math.gcd(mod,n);need=-i-t;assert need%g==0;k=(need//g*pow(mod//g,-1,n//g))%(n//g);t+=mod*k;mod=math.lcm(mod,n);t%=mod
save(13,wait*bus,t)
# Direct address enumeration (inputs have bounded floating bits).
m1={};m2={};largest=0
for l in text(14).splitlines():
 if l.startswith('mask'):
  mask=l.split(' = ')[1];floating=[35-i for i,c in enumerate(mask) if c=='X'];largest=max(largest,len(floating));ones=int(mask.replace('X','0'),2);keep=int(mask.replace('X','1'),2)
 else:
  address,value=map(int,re.findall(r'\d+',l));m1[address]=(value|ones)&keep;addresses=[address|ones]
  for bit in floating:addresses=[a&~(1<<bit) for a in addresses]+[a|(1<<bit) for a in addresses]
  for address in addresses:m2[address]=value
save(14,sum(m1.values()),sum(m2.values()),{'max_floating_bits':largest})
# Memory game: bounded typed last-seen array, no full spoken history.
from array import array
seed=[0,8,15,2,12,1,4];limit=30000000;seen=array('i',[0])*limit
for turn,v in enumerate(seed[:-1],1):seen[v]=turn
last=seed[-1];p1=None
for turn in range(len(seed),limit):
 prev=seen[last];seen[last]=turn;last=turn-prev if prev else 0
 if turn+1==2020:p1=last
save(15,p1,last,{'seed':seed,'input_path':'F# Day15.fs inline seed','input_sha256':None});del seen
# Ticket constraints and bipartite assignment through recursive matching.
rules_s,own_s,near_s=text(16).split('\n\n');rules={}
for l in rules_s.splitlines():
 name,body=l.split(': ');rules[name]=list(map(int,re.findall(r'\d+',body)))
def accepts(name,x):
 a,b,c,d=rules[name];return a<=x<=b or c<=x<=d
own=list(map(int,own_s.splitlines()[1].split(',')));near=[list(map(int,l.split(','))) for l in near_s.splitlines()[1:]];valid=[own];error=0
for ticket in near:
 invalid=[x for x in ticket if not any(accepts(k,x) for k in rules)];error+=sum(invalid)
 if not invalid:valid.append(ticket)
candidates={k:{i for i in range(len(own)) if all(accepts(k,row[i]) for row in valid)} for k in rules}
solutions=[]
def match(todo,used,result):
 if not todo:solutions.append(result);return
 k=min(todo,key=lambda k:len(candidates[k]-used))
 for pos in candidates[k]-used:match(todo-{k},used|{pos},dict(result,**{k:pos}))
match(set(rules),set(),{});assert len(solutions)==1
prod=math.prod(own[pos] for k,pos in solutions[0].items() if k.startswith('departure'));save(16,error,prod,{'assignment':solutions[0]})
# Sparse cellular automata via neighbour frequency, dimension parameter.
def life(dim):
 active={(x,y)+(0,)*(dim-2) for y,row in enumerate(text(17).splitlines()) for x,c in enumerate(row) if c=='#'};steps=[d for d in itertools.product((-1,0,1),repeat=dim) if any(d)]
 for _ in range(6):
  counts=collections.Counter(tuple(x+dx for x,dx in zip(p,d)) for p in active for d in steps);active={p for p,n in counts.items() if n==3 or n==2 and p in active}
 return len(active)
save(17,life(3),life(4))
# Arithmetic uses shunting-yard operator stacks (no Go parser dependency).
def evaluate(expr,addition_first):
 prec={'+':2 if addition_first else 1,'*':1};nums=[];ops=[]
 def apply():
  b=nums.pop();a=nums.pop();op=ops.pop();nums.append(a+b if op=='+' else a*b)
 for token in re.findall(r'\d+|[()+*]',expr):
  if token.isdigit():nums.append(int(token))
  elif token=='(':ops.append(token)
  elif token==')':
   while ops[-1]!='(':apply()
   ops.pop()
  else:
   while ops and ops[-1]!='(' and prec[ops[-1]]>=prec[token]:apply()
   ops.append(token)
 while ops:apply()
 assert len(nums)==1;return nums[0]
save(18,sum(evaluate(l,False) for l in text(18).splitlines()),sum(evaluate(l,True) for l in text(18).splitlines()))
# Grammar matcher: sets of consumed offsets, memoised per message.
rules_s,messages=text(19).split('\n\n');rules={}
for l in rules_s.splitlines():
 k,body=l.split(': ');rules[int(k)]=body.strip('"') if '"' in body else [list(map(int,s.split())) for s in body.split('|')]
def grammar_count(rulemap):
 total=0
 for message in messages.splitlines():
  @functools.cache
  def ends(rule,start):
   if start>=len(message):return frozenset()
   body=rulemap[rule]
   if isinstance(body,str):return frozenset([start+1]) if message[start]==body else frozenset()
   out=set()
   for sequence in body:
    positions={start}
    for child in sequence:positions={end for pos in positions for end in ends(child,pos)}
    out|=positions
   return frozenset(out)
  total+=len(message) in ends(0,0)
 return total
p1=grammar_count(rules);rules[8]=[[42],[42,8]];rules[11]=[[42,31],[42,11,31]];save(19,p1,grammar_count(rules))

print(json.dumps(OUT,indent=2))
