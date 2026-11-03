from pathlib import Path
import argparse
import hashlib,json,re,itertools,collections,functools,math,operator,time
parser = argparse.ArgumentParser(description='Recompute the independent 2017 answers.')
parser.add_argument('--inputs', type=Path, default=Path(__file__).resolve().parents[2] / 'inputs' / '2017')
parser.add_argument('--output', type=Path, required=True)
args = parser.parse_args()
ROOT = args.inputs
OUT = {}
INLINE={3:'265149',14:'ffayrhll',17:'376'}
SOURCE={1:(1144,1194),2:(30994,233),3:(438,266330),4:(337,231),5:(394829,31150702),6:(11137,1037),7:('vtzay',910),8:(5849,6702),9:(7616,3838),10:(52070,'7f94112db4e32e19cf6502073c66f9bb'),11:(705,1469),12:(169,179),13:(2164,3861798),14:(8190,1134),15:(594,328),16:('bijankplfgmeodhc','bpjahknliomefdgc'),17:(777,39289581),18:(3188,7112),19:('LOHMDQATP',16492),20:(243,648),21:(179,2766750),22:(5570,2512022),23:(4225,905),24:(1695,1673),25:(3578,None)}
def raw(d):
 text = (ROOT/f'{d:02}.txt').read_text()
 if d in INLINE:
  assert text.strip() == INLINE[d], f'inline input mismatch for day {d}'
  return INLINE[d]
 return text
def save(d,a,b=None,**extra):
 OUT[str(d)]={'part1':str(a),'part2':str(b) if b is not None else None,'source_expected':[str(x) if x is not None else None for x in SOURCE[d]],'input_path':f'inputs/2017/{d:02}.txt','inline_input':INLINE.get(d),'source_input_sha256':hashlib.sha256(raw(d).encode()).hexdigest(),'fixture_sha256':hashlib.sha256((ROOT/f'{d:02}.txt').read_bytes()).hexdigest(),'provenance':'independent Python oracle; no Go outputs used',**extra};print(d,a,b,flush=True)
# Fast/sequential batch; CPU-heavy15/17/22/25 run later explicitly.
s=raw(1).strip();save(1,sum(int(c) for i,c in enumerate(s) if c==s[(i+1)%len(s)]),sum(int(c) for i,c in enumerate(s) if c==s[(i+len(s)//2)%len(s)]))
rows=[list(map(int,l.split())) for l in raw(2).splitlines()];quot=[]
for row in rows:
 matches=[max(a,b)//min(a,b) for a,b in itertools.combinations(row,2) if max(a,b)%min(a,b)==0];assert len(matches)==1;quot+=matches
save(2,sum(max(row)-min(row) for row in rows),sum(quot))
# Literal spiral positions and sparse8-neighbour accumulated values.
target=int(raw(3));p=0j;v=1j;direction=1;grid={0j:1};index=1;stress=None;distance=None
for length in itertools.count(1):
 for _ in range(2):
  for _ in range(length):
   p+=direction;index+=1;value=sum(grid.get(p+dx+dy*1j,0) for dx in (-1,0,1) for dy in (-1,0,1) if dx or dy);grid[p]=value
   if value>target and stress is None:stress=value
   if index==target:distance=int(abs(p.real)+abs(p.imag))
  direction*=1j
 if stress is not None and distance is not None:break
save(3,distance,stress)
words=[l.split() for l in raw(4).splitlines()];save(4,sum(len(w)==len(set(w)) for w in words),sum(len(w)==len({''.join(sorted(x)) for x in w}) for w in words))
# Full jump execution (31Msteps), no inferred skipping.
values=list(map(int,raw(5).splitlines()))
def jumps(part2):
 a=values[:];ip=steps=0
 while 0<=ip<len(a):
  off=a[ip];a[ip]+=(-1 if part2 and off>=3 else 1);ip+=off;steps+=1
 return steps
save(5,jumps(False),jumps(True))
banks=list(map(int,raw(6).split()));seen={};steps=0
while tuple(banks) not in seen:
 seen[tuple(banks)]=steps;i=max(range(len(banks)),key=lambda i:banks[i]);n=banks[i];banks[i]=0;q,rem=divmod(n,len(banks));banks=[v+q for v in banks]
 for j in range(rem):banks[(i+j+1)%len(banks)]+=1
 steps+=1
save(6,steps,steps-seen[tuple(banks)])
children={};weights={}
for l in raw(7).splitlines():
 m=re.fullmatch(r'(\w+) \((\d+)\)(?: -> (.*))?',l);k,w,cs=m.groups();weights[k]=int(w);children[k]=cs.split(', ') if cs else []
root=(set(weights)-{c for cs in children.values() for c in cs}).pop();correction=[]
def tower(k):
 totals=[tower(c) for c in children[k]];freq=collections.Counter(totals)
 if len(freq)>1 and not correction:
  odd=next(v for v,n in freq.items() if n==1);want=next(v for v,n in freq.items() if n>1);wrong=children[k][totals.index(odd)];correction.append(weights[wrong]+want-odd)
 return weights[k]+sum(totals)
tower(root);save(7,root,correction[0])
registers=collections.defaultdict(int);ever=0;cmp={'<':operator.lt,'>':operator.gt,'<=':operator.le,'>=':operator.ge,'==':operator.eq,'!=':operator.ne}
for l in raw(8).splitlines():
 dest,op,num,_,cond,rel,val=l.split();registers[dest]+=0
 if cmp[rel](registers[cond],int(val)):registers[dest]+=int(num)*(1 if op=='inc' else -1);ever=max(ever,registers[dest])
save(8,max(registers.values()),ever)
garbage=False;cancel=False;depth=score=count=0
for c in raw(9).strip():
 if cancel:cancel=False;continue
 if c=='!':cancel=True;continue
 if garbage:
  if c=='>':garbage=False
  else:count+=1
 elif c=='<':garbage=True
 elif c=='{':depth+=1;score+=depth
 elif c=='}':depth-=1
save(9,score,count)
def sparse(lengths,rounds,n=256):
 a=list(range(n));pos=skip=0
 for _ in range(rounds):
  for size in lengths:
   for j in range(size//2):x=(pos+j)%n;y=(pos+size-j-1)%n;a[x],a[y]=a[y],a[x]
   pos=(pos+size+skip)%n;skip+=1
 return a

def knot(s):
 a=sparse(list(s.encode())+[17,31,73,47,23],64);return bytes(functools.reduce(operator.xor,a[i:i+16]) for i in range(0,256,16))
s=raw(10).strip();a=sparse(list(map(int,s.split(','))),1);save(10,a[0]*a[1],knot(s).hex())
dirs={'n':(0,1,-1),'ne':(1,0,-1),'se':(1,-1,0),'s':(0,-1,1),'sw':(-1,0,1),'nw':(-1,1,0)};p=[0,0,0];furthest=0
for c in raw(11).strip().split(','):
 p=[a+b for a,b in zip(p,dirs[c])];furthest=max(furthest,max(map(abs,p)))
save(11,max(map(abs,p)),furthest)
graph={}
for l in raw(12).splitlines():
 v,cs=l.split(' <-> ');graph[int(v)]=set(map(int,cs.split(', ')))
def component(start):
 seen={start};stack=[start]
 while stack:
  v=stack.pop()
  for c in graph[v]-seen:seen.add(c);stack.append(c)
 return seen
first=component(0);left=set(graph);groups=0
while left:left-=component(next(iter(left)));groups+=1
save(12,len(first),groups)
firewall=[tuple(map(int,l.split(': '))) for l in raw(13).splitlines()];severity=sum(d*r for d,r in firewall if r==1 or d%(2*(r-1))==0);periods=[(d,2*(r-1)) for d,r in firewall];delay=next(t for t in itertools.count() if all((t+d)%p for d,p in periods));save(13,severity,delay)
used=set()
for y in range(128):
 bits=int.from_bytes(knot(raw(14)+'-'+str(y)),'big')
 for x in range(128):
  if bits&(1<<(127-x)):used.add((x,y))
a=len(used);regions=0
while used:
 start=used.pop();stack=[start];regions+=1
 while stack:
  x,y=stack.pop()
  for n in ((x-1,y),(x+1,y),(x,y-1),(x,y+1)):
   if n in used:used.remove(n);stack.append(n)
save(14,a,regions)
ops=raw(16).strip().split(',')
def dance(s):
 s=list(s)
 for op in ops:
  if op[0]=='s':n=int(op[1:])%len(s);s=s[-n:]+s[:-n] if n else s
  else:
   a,b=op[1:].split('/');a,b=(int(a),int(b)) if op[0]=='x' else (s.index(a),s.index(b));s[a],s[b]=s[b],s[a]
 return ''.join(s)
s='abcdefghijklmnop';cycle=[]
while s not in cycle:cycle.append(s);s=dance(s)
assert s==cycle[0];save(16,dance(cycle[0]),cycle[1000000000%len(cycle)],cycle_length=len(cycle))
# Duet plain and two-process deterministic fair scheduler.
program=[l.split() for l in raw(18).splitlines()]
def sound():
 r=collections.defaultdict(int);ip=0;last=None;steps=0
 def val(s):return r[s] if s.isalpha() else int(s)
 while 0<=ip<len(program):
  steps+=1;assert steps<10000000;op,x,*ys=program[ip];y=ys[0] if ys else None
  if op=='snd':last=val(x)
  elif op=='rcv':
   if val(x):return last
  elif op=='set':r[x]=val(y)
  elif op=='add':r[x]+=val(y)
  elif op=='mul':r[x]*=val(y)
  elif op=='mod':r[x]%=val(y)
  elif op=='jgz' and val(x)>0:ip+=val(y);continue
  ip+=1
 raise ValueError('terminated before recovery')
def duet():
 rs=[collections.defaultdict(int,p=i) for i in range(2)];ips=[0,0];queues=[collections.deque(),collections.deque()];sent=[0,0];steps=0
 while True:
  progress=False
  for i in (0,1):
   ip=ips[i]
   if not 0<=ip<len(program):continue
   r=rs[i];val=lambda s:r[s] if s.isalpha() else int(s);op,x,*ys=program[ip];y=ys[0] if ys else None
   if op=='rcv' and not queues[i]:continue
   progress=True;steps+=1;assert steps<20000000
   if op=='snd':queues[1-i].append(val(x));sent[i]+=1
   elif op=='rcv':r[x]=queues[i].popleft()
   elif op=='set':r[x]=val(y)
   elif op=='add':r[x]+=val(y)
   elif op=='mul':r[x]*=val(y)
   elif op=='mod':r[x]%=val(y)
   elif op=='jgz' and val(x)>0:ips[i]+=val(y);continue
   ips[i]+=1
  if not progress:return sent[1]
save(18,sound(),duet())
grid=raw(19).splitlines();x=grid[0].index('|');y=0;dx,dy=0,1;letters='';steps=0
def at(x,y):return grid[y][x] if 0<=y<len(grid) and 0<=x<len(grid[y]) else ' '
while at(x,y)!=' ':
 c=at(x,y);steps+=1
 if c.isalpha():letters+=c
 if c=='+':
  choices=[(a,b) for a,b in ((1,0),(-1,0)) if at(x+a,y+b)!=' '] if dx==0 else [(a,b) for a,b in ((0,1),(0,-1)) if at(x+a,y+b)!=' '];assert len(choices)==1;dx,dy=choices[0]
 x+=dx;y+=dy
save(19,letters,steps)
# Exact eventual polynomial and all integer collision events including t0.
particles=[tuple(map(int,re.findall(r'-?\d+',l))) for l in raw(20).splitlines()]
def eventual(p):
 coef=[0,0,0]
 for axis in range(3):
  pp,v,a=p[axis],p[axis+3],p[axis+6];sgn=1 if next((x for x in (a,v,pp) if x),0)>0 else -1
  for k,value in enumerate((a,a+2*v,2*pp)):coef[k]+=sgn*value
 return tuple(coef)
a=min(range(len(particles)),key=lambda i:(eventual(particles[i]),i))
def roots(a,b,c):
 if a==0:
  if b==0:return None if c==0 else set()
  return {-c//b} if (-c)%b==0 and -c//b>=0 else set()
 disc=b*b-4*a*c
 if disc<0:return set()
 s=math.isqrt(disc)
 if s*s!=disc:return set()
 return {n//(2*a) for n in (-b+s,-b-s) if n%(2*a)==0 and n//(2*a)>=0}
events=collections.defaultdict(list)
for i,p in enumerate(particles):
 for j in range(i):
  q=particles[j];times=None
  for ax in range(3):
   rr=roots(p[ax+6]-q[ax+6],p[ax+6]-q[ax+6]+2*(p[ax+3]-q[ax+3]),2*(p[ax]-q[ax]))
   if rr is not None:times=rr if times is None else times&rr
   if times==set():break
  if times is None:times={0}
  for t in times:events[t].append((i,j))
alive=set(range(len(particles)))
for t,pairs in sorted(events.items()):
 removed={i for a,b in pairs if a in alive and b in alive for i in (a,b)};alive-=removed
save(20,a,len(alive),event_times=sorted(events),initial_duplicate_pairs=len(events.get(0,[])),eventual_coefficients=eventual(particles[a]))
# Fractal full dense image, orientation lookup generated independently.
def rotate(pattern):return tuple(''.join(row[i] for row in pattern[::-1]) for i in range(len(pattern)))
rules={}
for line in raw(21).splitlines():
 left,right=line.split(' => ');p=tuple(left.split('/'));out=tuple(right.split('/'))
 for _ in range(4):
  for variant in (p,tuple(row[::-1] for row in p)):
   if variant in rules:assert rules[variant]==out
   rules[variant]=out
  p=rotate(p)
image=('.#.','..#','###');first=None
for step in range(1,19):
 block=2 if len(image)%2==0 else 3;size=len(image)//block;out=[]
 for gy in range(size):
  tiles=[rules[tuple(row[gx*block:(gx+1)*block] for row in image[gy*block:(gy+1)*block])] for gx in range(size)]
  out.extend(''.join(tile[y] for tile in tiles) for y in range(block+1))
 image=tuple(out)
 if step==5:first=sum(row.count('#') for row in image)
save(21,first,sum(row.count('#') for row in image))
# Actual part1 VM and part2 independently proved composite progression.
program=[l.split() for l in raw(23).splitlines()];r=collections.defaultdict(int);ip=mul=0
val=lambda x:r[x] if x.isalpha() else int(x)
while 0<=ip<len(program):
 op,x,y=program[ip]
 if op=='set':r[x]=val(y)
 elif op=='sub':r[x]-=val(y)
 elif op=='mul':r[x]*=val(y);mul+=1
 elif op=='jnz' and val(x):ip+=val(y);continue
 ip+=1
r=collections.defaultdict(int,a=1);ip=0
while ip!=8:
 op,x,y=program[ip]
 if op=='set':r[x]=val(y)
 elif op=='sub':r[x]-=val(y)
 elif op=='mul':r[x]*=val(y)
 elif op=='jnz' and val(x):ip+=val(y);continue
 ip+=1
start,end=r['b'],r['c'];stride=-int(program[30][2]);composites=[n for n in range(start,end+1,stride) if any(n%d==0 for d in range(2,math.isqrt(n)+1))];save(23,mul,len(composites),start=start,end=end,stride=stride,progression_length=len(range(start,end+1,stride)),composite_values=composites)
# Index-distinct bridge components and exhaustive remaining-mask recursion.
components=[tuple(map(int,l.split('/'))) for l in raw(24).splitlines()]
@functools.cache
def bridge(port,mask):
 best=0;longest=(0,0)
 for i,(a,b) in enumerate(components):
  if mask>>i&1 and port in (a,b):
   strength,ln,ls=bridge(b if port==a else a,mask^(1<<i));best=max(best,a+b+strength);longest=max(longest,(ln+1,ls+a+b))
 return best,*longest
best,ln,ls=bridge(0,(1<<len(components))-1);save(24,best,ls,components=len(components),distinct_components=len(set(components)))
# Explicit serial heavy continuation with wall deadlines.
def heavy_generators():
 started=time.monotonic();a,b=map(int,re.findall(r'\d+',raw(15)));aa,bb=a,b;matches=0
 for i in range(40000000):
  a=a*16807%2147483647;b=b*48271%2147483647;matches+=(a&65535)==(b&65535)
  if i%1000000==0 and time.monotonic()-started>70:raise TimeoutError('generator oracle budget')
 first=matches;a,b=aa,bb;matches=0
 for i in range(5000000):
  while True:
   a=a*16807%2147483647
   if a%4==0:break
  while True:
   b=b*48271%2147483647
   if b%8==0:break
  matches+=(a&65535)==(b&65535)
  if i%100000==0 and time.monotonic()-started>70:raise TimeoutError('generator oracle budget')
 save(15,first,matches,seeds=[aa,bb])

def heavy_spinlock():
 started=time.monotonic();step=int(raw(17));buffer=[0];p=0
 for n in range(1,2018):p=(p+step)%len(buffer)+1;buffer.insert(p,n)
 first=buffer[(p+1)%len(buffer)];p=after=0
 for n in range(1,50000001):
  p=(p+step)%n+1
  if p==1:after=n
  if n%1000000==0 and time.monotonic()-started>30:raise TimeoutError('spinlock oracle budget')
 save(17,first,after)

def heavy_virus():
 started=time.monotonic();rows=raw(22).splitlines();initial={(x-len(row)//2,y-len(rows)//2):2 for y,row in enumerate(rows) for x,c in enumerate(row) if c=='#'}
 def infect(evolved,bursts):
  states=dict(initial);x=y=d=infected=0;directions=[(0,-1),(1,0),(0,1),(-1,0)]
  for i in range(bursts):
   old=states.get((x,y),0)
   if evolved:d=(d+(-1,0,1,2)[old])%4;new=(old+1)%4
   else:d=(d+(1 if old==2 else -1))%4;new=0 if old==2 else 2
   if new:states[x,y]=new
   else:states.pop((x,y),None)
   infected+=new==2;dx,dy=directions[d];x+=dx;y+=dy
   if i%1000000==0 and time.monotonic()-started>50:raise TimeoutError('virus oracle budget')
  return infected
 save(22,infect(False,10000),infect(True,10000000))

def heavy_turing():
 started=time.monotonic();blueprint=raw(25);state=re.search(r'Begin in state ([A-Z])\.',blueprint)[1];steps=int(re.search(r'after (\d+) steps',blueprint)[1]);rules={}
 for name,block in re.findall(r'In state ([A-Z]):(.*?)(?=\nIn state |\Z)',blueprint,re.S):
  transitions=re.findall(r'If the current value is ([01]):\s+- Write the value ([01])\.\s+- Move one slot to the (left|right)\.\s+- Continue with state ([A-Z])\.',block)
  assert len(transitions)==2
  for bit,write,move,nxt in transitions:rules[name,int(bit)]=(int(write),1 if move=='right' else -1,nxt)
 tape=set();cursor=0
 for i in range(steps):
  write,move,state=rules[state,cursor in tape]
  if write:tape.add(cursor)
  else:tape.discard(cursor)
  cursor+=move
  if i%1000000==0 and time.monotonic()-started>40:raise TimeoutError('Turing oracle budget')
 save(25,len(tape),None,initial_state=re.search(r'Begin in state ([A-Z])\.',blueprint)[1],steps=steps,states=len(rules)//2)

heavy_generators()
heavy_spinlock()
heavy_virus()
heavy_turing()
assert len(OUT) == 25
for day, record in OUT.items():
    assert [record['part1'], record['part2']] == record['source_expected'], day
args.output.write_text(json.dumps(dict(sorted(OUT.items(), key=lambda pair: int(pair[0]))), indent=2) + '\n')
