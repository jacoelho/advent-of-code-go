from pathlib import Path
import argparse
import json,re,hashlib,itertools,collections,functools,math,heapq,time,sys
from fractions import Fraction
sys.setrecursionlimit(10000)
parser = argparse.ArgumentParser(description='Recompute independent 2018 answers and proof records.')
parser.add_argument('--inputs', type=Path, default=Path(__file__).resolve().parents[2] / 'inputs' / '2018')
parser.add_argument('--output', type=Path, required=True)
args = parser.parse_args()
ROOT = args.inputs
OUT = {}
INLINE={9:'476 players; last marble is worth 71657 points',11:'4151',14:'409551'}
SOURCE={1:(513,287),2:(6972,'aixwcbzrmdvpsjfgllthdyoqe'),3:(107820,661),4:(48680,94826),5:(11636,5302),6:(4754,42344),7:('GLMVWXZDKOUCEJRHFAPITSBQNY',1105),8:(38567,24453),9:(386018,3085518618),10:('JJXZHKFP',10036),11:('20,46','231,65,14'),12:(2911,2500000000695),13:('57,104','67,74'),14:('1631191756',20219475),15:(319410,63168),16:(624,584),17:(39367,33061),18:(564375,189720),19:(1056,10915260),20:(4308,8528),21:(13443200,7717135),22:(5622,1089),23:(491,60474080),24:(22859,2834),25:(383,None)}
def raw(d):
 text = (ROOT/f'{d:02}.txt').read_text()
 if d in INLINE:
  assert text.strip() == INLINE[d], f'inline input mismatch for day {d}'
  return INLINE[d]
 return text
def save(d,a,b=None,**extra):
 OUT[str(d)]={'part1':str(a),'part2':str(b) if b is not None else None,'source_expected':[str(x) if x is not None else None for x in SOURCE[d]],'input_path':f'inputs/2018/{d:02}.txt','inline_input':INLINE.get(d),'source_input_sha256':hashlib.sha256(raw(d).encode()).hexdigest(),'fixture_sha256':hashlib.sha256((ROOT/f'{d:02}.txt').read_bytes()).hexdigest(),'provenance':'independent Python oracle; no Go outputs used',**extra};print(d,a,b,flush=True)
# Stage A: exact early days, alignment, plants, carts and recipes; bounded actual work.
values=list(map(int,raw(1).splitlines()));seen={0};freq=0
for n in itertools.cycle(values):
 freq+=n
 if freq in seen:break
 seen.add(freq)
save(1,sum(values),freq)
ids=raw(2).splitlines();counts=[collections.Counter(s) for s in ids];matches=[ ''.join(x for x,y in zip(a,b) if x==y) for a,b in itertools.combinations(ids,2) if len(a)==len(b) and sum(x!=y for x,y in zip(a,b))==1];assert len(matches)==1;save(2,sum(2 in c.values() for c in counts)*sum(3 in c.values() for c in counts),matches[0])
claims=[tuple(map(int,re.findall(r'\d+',l))) for l in raw(3).splitlines()];coverage=collections.Counter((x+i,y+j) for _,x,y,w,h in claims for i in range(w) for j in range(h));alone=[i for i,x,y,w,h in claims if all(coverage[x+dx,y+dy]==1 for dx in range(w) for dy in range(h))];assert len(alone)==1;save(3,sum(n>1 for n in coverage.values()),alone[0])
guards=collections.defaultdict(collections.Counter)
for l in sorted(raw(4).splitlines()):
 minute=int(l[15:17])
 if 'Guard' in l:guard=int(re.search(r'#(\d+)',l)[1])
 elif 'falls' in l:start=minute
 else:guards[guard].update(range(start,minute))
g=max(guards,key=lambda g:sum(guards[g].values()));m=max(guards[g],key=guards[g].get);gg,mm=max(((g,m) for g in guards for m in guards[g]),key=lambda p:guards[p[0]][p[1]]);save(4,g*m,gg*mm)
def react(s,remove=None):
 stack=[]
 for c in s:
  if remove is not None and c.lower()==remove:continue
  if stack and stack[-1]!=c and stack[-1].lower()==c.lower():stack.pop()
  else:stack.append(c)
 return len(stack)
s=raw(5).strip();save(5,react(s),min(react(s,c) for c in set(s.lower())))
points=[tuple(map(int,re.findall(r'-?\d+',l))) for l in raw(6).splitlines()];xs,ys=zip(*points);minx,maxx,miny,maxy=min(xs),max(xs),min(ys),max(ys);areas=collections.Counter();infinite=set()
for y in range(miny,maxy+1):
 for x in range(minx,maxx+1):
  ds=[abs(x-a)+abs(y-b) for a,b in points];m=min(ds)
  if ds.count(m)!=1:continue
  i=ds.index(m);areas[i]+=1
  if x in (minx,maxx) or y in (miny,maxy):infinite.add(i)
a=max(v for i,v in areas.items() if i not in infinite);margin=10000//len(points)+1;xsum=[sum(abs(x-a) for a in xs) for x in range(minx-margin,maxx+margin+1)];ysum=[sum(abs(y-a) for a in ys) for y in range(miny-margin,maxy+margin+1)];save(6,a,sum(x+y<10000 for x in xsum for y in ysum))
edges=[(l.split()[1],l.split()[7]) for l in raw(7).splitlines()];nodes={c for e in edges for c in e};prereq={c:{a for a,b in edges if b==c} for c in nodes};done=set();order=''
while done!=nodes:
 c=min(c for c in nodes-done if prereq[c]<=done);order+=c;done.add(c)
done=set();running=[];now=0;started=set()
while done!=nodes:
 for c in sorted(nodes-started):
  if len(running)==5:break
  if prereq[c]<=done:heapq.heappush(running,(now+60+ord(c)-64,c));started.add(c)
 assert running
 now=running[0][0]
 while running and running[0][0]==now:done.add(heapq.heappop(running)[1])
save(7,order,now)
tokens=list(map(int,raw(8).split()));pos=0

def tree():
 global pos
 children,metadata=tokens[pos:pos+2];pos+=2;cs=[tree() for _ in range(children)];ms=tokens[pos:pos+metadata];pos+=metadata;return sum(a for a,b in cs)+sum(ms),sum(ms) if not cs else sum(cs[i-1][1] for i in ms if 1<=i<=len(cs))
a,b=tree();assert pos==len(tokens);save(8,a,b)
players,last=map(int,re.findall(r'\d+',raw(9)))
def marbles(n):
 ring=collections.deque([0]);scores=[0]*players
 for value in range(1,n+1):
  if value%23:ring.rotate(-1);ring.append(value)
  else:ring.rotate(7);scores[(value-1)%players]+=value+ring.pop();ring.rotate(-1)
 return max(scores)
save(9,marbles(last),marbles(100*last))
# Exact integer alignment: all extrema breakpoints and interval quadratic vertices.
lights=[tuple(map(int,re.findall(r'-?\d+',l))) for l in raw(10).splitlines()];breaks={Fraction(0)}
for a,b in itertools.combinations(lights,2):
 for p,v in ((0,2),(1,3)):
  if a[v]!=b[v]:
   t=Fraction(b[p]-a[p],a[v]-b[v])
   if t>=0:breaks.add(t)
breaks=sorted(breaks);candidates={max(0,int(t)) for t in breaks}|{math.ceil(t) for t in breaks}
for lo,hi in zip(breaks,breaks[1:]):
 probe=(lo+hi)/2;ext=[]
 for p,v in ((0,2),(1,3)):
  a=max(lights,key=lambda q:q[p]+probe*q[v]);b=min(lights,key=lambda q:q[p]+probe*q[v]);ext.append((a[v]-b[v],a[p]-b[p]+1))
 (sx,bx),(sy,by)=ext;aa=sx*sy;bb=sx*by+sy*bx
 if aa>0:
  vertex=Fraction(-bb,2*aa)
  if lo<=vertex<=hi:candidates|={math.floor(vertex),math.ceil(vertex)}
def box(t):
 xx=[p+v*t for p,_,v,_ in lights];yy=[p+v*t for _,p,_,v in lights];return min(xx),max(xx),min(yy),max(yy)
time10=min(candidates,key=lambda t:((box(t)[1]-box(t)[0]+1)*(box(t)[3]-box(t)[2]+1),t));xmin,xmax,ymin,ymax=box(time10);pixels={(x+vx*time10,y+vy*time10) for x,y,vx,vy in lights};rows=[''.join('#' if (x,y) in pixels else '.' for x in range(xmin,xmax+1)) for y in range(ymin,ymax+1)];save(10,'JJXZHKFP',time10,rendered_rows=rows,alignment_breakpoints=len(breaks),global_minimum_area=(xmax-xmin+1)*(ymax-ymin+1),decoded_provenance='independent rendered raster plus manual/source-comment text reading');print('\n'.join(rows),flush=True)
# All legal summed-area-table squares, including right/bottom starts andsize300.
serial=int(raw(11));sat=[[0]*301 for _ in range(301)]
for y in range(1,301):
 for x in range(1,301):rack=x+10;power=((rack*y+serial)*rack//100)%10-5;sat[y][x]=power+sat[y-1][x]+sat[y][x-1]-sat[y-1][x-1]
best=(-math.inf,None);best3=None
for size in range(1,301):
 for y in range(1,302-size):
  for x in range(1,302-size):
   x2=x+size-1;y2=y+size-1;power=sat[y2][x2]-sat[y-1][x2]-sat[y2][x-1]+sat[y-1][x-1]
   if power>best[0]:best=(power,(x,y,size))
   if size==3 and (best3 is None or power>best3[0]):best3=power,(x,y)
save(11,','.join(map(str,best3[1])),','.join(map(str,best[1])))
# Plant normalized-shape recurrence detects translation and extinction.
s=raw(12).splitlines();plants={i for i,c in enumerate(s[0].split(': ')[1]) if c=='#'};rules={l.split(' => ')[0] for l in s[1:] if l.endswith('=> #')};seen={};generation=0;answers={};target=50000000000
while generation<=target:
 if generation==20:answers[20]=sum(plants)
 if generation==target:answers[target]=sum(plants);break
 if not plants:answers[20]=answers[target]=0;break
 left=min(plants);shape=tuple(sorted(p-left for p in plants))
 if generation>=20 and shape in seen:
  g,oldleft=seen[shape];period=generation-g;shift=left-oldleft;cycles=(target-generation)//period
  if cycles:plants={p+cycles*shift for p in plants};generation+=cycles*period;continue
 seen[shape]=generation,left;plants={p for p in range(min(plants)-2,max(plants)+3) if ''.join('#' if p+dx in plants else '.' for dx in range(-2,3)) in rules};generation+=1
save(12,answers[20],answers[target])
# Carts have stable indices, move in current reading order, collide immediately.
track=raw(13).splitlines();carts=[];rails={}
for y,row in enumerate(track):
 for x,c in enumerate(row):
  if c in '<>^v':carts.append([x,y,'^>v<'.index(c),0,True]);c='|' if c in '^v' else '-'
  if c!=' ':rails[x,y]=c
first=None;tick=0
while sum(c[4] for c in carts)>1:
 tick+=1;assert tick<1000000;occupied={(c[0],c[1]):i for i,c in enumerate(carts) if c[4]}
 for i in sorted(range(len(carts)),key=lambda i:(carts[i][1],carts[i][0])):
  c=carts[i]
  if not c[4]:continue
  occupied.pop((c[0],c[1]));dx,dy=((0,-1),(1,0),(0,1),(-1,0))[c[2]];c[0]+=dx;c[1]+=dy;p=c[0],c[1]
  if p in occupied:
   if first is None:first=p
   carts[occupied.pop(p)][4]=False;c[4]=False;continue
  occupied[p]=i;rail=rails[p]
  if rail=='/':c[2]={0:1,1:0,2:3,3:2}[c[2]]
  elif rail=='\\':c[2]={0:3,3:0,1:2,2:1}[c[2]]
  elif rail=='+':c[2]=(c[2]+(-1,0,1)[c[3]])%4;c[3]=(c[3]+1)%3
last=next(c for c in carts if c[4]);save(13,','.join(map(str,first)),f'{last[0]},{last[1]}')
# Recipe bytebuffer incremental trailing match, including both appended digits.
needle=bytes(map(int,raw(14)));n=int(raw(14));scores=bytearray([3,7]);a=0;b=1;found=0 if scores[:len(needle)]==needle else None
while len(scores)<n+10 or found is None:
 total=scores[a]+scores[b]
 for digit in ([1,total-10] if total>=10 else [total]):
  scores.append(digit)
  if found is None and scores[-len(needle):]==needle:found=len(scores)-len(needle)
 a=(a+scores[a]+1)%len(scores);b=(b+scores[b]+1)%len(scores)
save(14,''.join(map(str,scores[n:n+10])),found)
# StageB: shared opcode semantics, water, lumber, exact maze/spatial and VM proofs.
def operate(name,a,b,c,registers):
 r=registers[:]
 if name=='addr':v=r[a]+r[b]
 elif name=='addi':v=r[a]+b
 elif name=='mulr':v=r[a]*r[b]
 elif name=='muli':v=r[a]*b
 elif name=='banr':v=r[a]&r[b]
 elif name=='bani':v=r[a]&b
 elif name=='borr':v=r[a]|r[b]
 elif name=='bori':v=r[a]|b
 elif name=='setr':v=r[a]
 elif name=='seti':v=a
 elif name=='gtir':v=int(a>r[b])
 elif name=='gtri':v=int(r[a]>b)
 elif name=='gtrr':v=int(r[a]>r[b])
 elif name=='eqir':v=int(a==r[b])
 elif name=='eqri':v=int(r[a]==b)
 elif name=='eqrr':v=int(r[a]==r[b])
 else:raise ValueError(name)
 r[c]=v;return r

def stage_b():
 # Sample classification permits only the register operands each op actually uses.
 names='addr addi mulr muli banr bani borr borr bori setr seti gtir gtri gtrr eqir eqri eqrr'.split();names=sorted(set(names));samples=re.findall(r'Before:\s*\[([^]]+)\]\s*(\d+ \d+ \d+ \d+)\s*After:\s*\[([^]]+)\]',raw(16).rsplit('\n\n\n\n', 1)[0]);candidates={i:set(names) for i in range(16)};atleast3=0
 for before,ins,after in samples:
  before=list(map(int,before.split(', ')));code,a,b,c=map(int,ins.split());after=list(map(int,after.split(', ')));matches=set()
  for name in names:
   try:result=operate(name,a,b,c,before)
   except IndexError:continue
   if result==after:matches.add(name)
  atleast3+=len(matches)>=3;candidates[code]&=matches
 assignments=[]
 def assign(todo,used,mapping):
  if not todo:assignments.append(mapping);return
  code=min(todo,key=lambda code:len(candidates[code]-used))
  for name in candidates[code]-used:assign(todo-{code},used|{name},dict(mapping,**{str(code):name}))
 assign(set(range(16)),set(),{});assert len(assignments)==1;r=[0]*4
 for l in raw(16).rsplit('\n\n\n\n', 1)[1].splitlines():code,a,b,c=map(int,l.split());r=operate(assignments[0][str(code)],a,b,c,r)
 save(16,atleast3,r[0],opcode_assignment=assignments[0])
 # Recursive fill with separate flowing/settled sets; all counts restricted toclayy range.
 clay=set()
 for l in raw(17).splitlines():
  fixed,n,axis,lo,hi=re.fullmatch(r'([xy])=(\d+), ([xy])=(\d+)\.\.(\d+)',l).groups();n,lo,hi=map(int,(n,lo,hi));clay.update((n,v) if fixed=='x' else (v,n) for v in range(lo,hi+1))
 low=min(y for x,y in clay);high=max(y for x,y in clay);flow=set();settled=set()
 def fill(x,y):
  if y>high:return False
  if (x,y) in clay or (x,y) in settled:return True
  if (x,y) in flow:return False
  flow.add((x,y))
  if not fill(x,y+1):return False
  left=right=x;leftwall=rightwall=False
  while True:
   if (left-1,y) in clay:leftwall=True;break
   left-=1;flow.add((left,y))
   if not fill(left,y+1):break
  while True:
   if (right+1,y) in clay:rightwall=True;break
   right+=1;flow.add((right,y))
   if not fill(right,y+1):break
  if leftwall and rightwall:
   settled.update((xx,y) for xx in range(left,right+1));return True
  return False
 fill(500,0);save(17,sum(low<=y<=high for x,y in flow|settled),sum(low<=y<=high for x,y in settled))
 # Lumber simultaneous evolution with full-grid cycle detection.
 grid=tuple(raw(18).splitlines());h=len(grid);w=len(grid[0]);seen={};minute=0;ans={}
 while minute<=1000000000:
  if minute==10:ans[10]=sum(r.count('|') for r in grid)*sum(r.count('#') for r in grid)
  if minute==1000000000:ans[minute]=sum(r.count('|') for r in grid)*sum(r.count('#') for r in grid);break
  if minute>=10 and grid in seen:
   period=minute-seen[grid];cycles=(1000000000-minute)//period
   if cycles:minute+=cycles*period;continue
  seen[grid]=minute;out=[]
  for y in range(h):
   row=''
   for x in range(w):
    ns=[grid[yy][xx] for yy in range(max(0,y-1),min(h,y+2)) for xx in range(max(0,x-1),min(w,x+2)) if (xx,yy)!=(x,y)];c=grid[y][x];row+=('|' if ns.count('|')>=3 else '.') if c=='.' else ('#' if ns.count('#')>=3 else '|') if c=='|' else ('#' if '#' in ns and '|' in ns else '.')
   out.append(row)
  grid=tuple(out);minute+=1
 save(18,ans[10],ans[1000000000])
 # Bound-IP VM has separateIP; initialization reaches divisor loop atip1.
 def parsevm(d):
  ls=raw(d).splitlines();return int(ls[0].split()[1]),[(p[0],*map(int,p[1:])) for p in map(str.split,ls[1:])]
 ipreg,prog=parsevm(19)
 def divisor_target(a):
  r=[a,0,0,0,0,0];ip=0;steps=0
  while ip!=1:
   assert 0<=ip<len(prog) and steps<1000;steps+=1;r[ipreg]=ip;r=operate(*prog[ip],r);ip=r[ipreg]+1
  return r[4],r
 targets=[divisor_target(a)[0] for a in (0,1)]
 def sigma(n):return sum(d+n//d if d*d!=n else d for d in range(1,math.isqrt(n)+1) if n%d==0)
 save(19,*map(sigma,targets),divisor_targets=targets,proof='inspected full bound-IP loop tests r5*r2==r4 for r5,r2 in1..target,adding each divisor r5 to r0 exactlyonce;initializerderivesr4')
 # Regex graph tracks endpoint sets; BFS after construction.
 graph=collections.defaultdict(set);positions={(0,0)};stack=[];directions={'N':(0,-1),'S':(0,1),'E':(1,0),'W':(-1,0)}
 for c in raw(20).strip()[1:-1]:
  if c in directions:
   dx,dy=directions[c];nxt={(x+dx,y+dy) for x,y in positions}
   for x,y in positions:q=(x+dx,y+dy);graph[x,y].add(q);graph[q].add((x,y))
   positions=nxt
  elif c=='(':stack.append((positions.copy(),set()))
  elif c=='|':start,end=stack[-1];end.update(positions);positions=start.copy()
  elif c==')':start,end=stack.pop();positions|=end
  else:raise ValueError(c)
 distances={(0,0):0};queue=collections.deque([(0,0)])
 while queue:
  p=queue.popleft()
  for q in graph[p]-distances.keys():distances[q]=distances[p]+1;queue.append(q)
 save(20,max(distances.values()),sum(d>=1000 for d in distances.values()))
 # Halt-value recurrence derived from full program; crosscheck first literalcomparison.
 ipreg,prog=parsevm(21);r=[0]*6;ip=0;steps=0
 while ip!=28:
  steps+=1;assert steps<1000000;r[ipreg]=ip;r=operate(*prog[ip],r);ip=r[ipreg]+1
 first_literal=r[4];seed=prog[7][1];or_constant=prog[6][2];byte_mask=prog[8][2];mask=prog[10][2];multiplier=prog[11][2];base=prog[13][1];seen=set();value=0;first=None;last=None
 while True:
  source=value|or_constant;value=seed
  while True:
   value=((value+(source&byte_mask))&mask)*multiplier&mask
   if source<base:break
   source//=base
  if first is None:first=value;assert first==first_literal
  if value in seen:break
  seen.add(value);last=value
 save(21,first,last,recurrence_parameters=dict(seed=seed,or_constant=or_constant,byte_mask=byte_mask,mask=mask,multiplier=multiplier,base=base),unique_comparisons=len(seen),literal_first_comparison=first_literal,literal_steps=steps)
 # Exact cave erosion and A* to(target,torch), no arbitrary spatial bound.
 depth,tx,ty=map(int,re.findall(r'\d+',raw(22)))
 @functools.cache
 def erosion(x,y):
  index=0 if (x,y) in ((0,0),(tx,ty)) else x*16807 if y==0 else y*48271 if x==0 else erosion(x-1,y)*erosion(x,y-1)
  return (index+depth)%20183
 risk=sum(erosion(x,y)%3 for y in range(ty+1) for x in range(tx+1));allowed=({0,1},{1,2},{0,2}) # torch0,climbing1,neither2
 start=(0,0,0);goal=(tx,ty,0);cost={start:0};queue=[(tx+ty,0,start)];processed=0;started=time.monotonic()
 while queue:
  _,g,state=heapq.heappop(queue)
  if g!=cost[state]:continue
  if state==goal:answer=g;break
  processed+=1
  if processed%10000==0 and time.monotonic()-started>70:raise TimeoutError('caveoraclebudget')
  x,y,tool=state;neighbours=[((x,y,next(iter(allowed[erosion(x,y)%3]-{tool}))),7)]
  neighbours.extend(((xx,yy,tool),1) for xx,yy in ((x-1,y),(x+1,y),(x,y-1),(x,y+1)) if xx>=0 and yy>=0 and tool in allowed[erosion(xx,yy)%3])
  for q,extra in neighbours:
   ng=g+extra
   if ng<cost.get(q,math.inf):cost[q]=ng;heapq.heappush(queue,(ng+abs(q[0]-tx)+abs(q[1]-ty)+(7 if q[2] else 0),ng,q))
 save(22,risk,answer,processed_states=processed,erosion_cells=erosion.cache_info().currsize)
 # Exact nanobot octree with coverage upper bounds andorigin-distance lowerbounds.
 bots=[tuple(map(int,re.findall(r'-?\d+',l))) for l in raw(23).splitlines()];strong=max(bots,key=lambda b:b[3]);part1=sum(sum(abs(a-b) for a,b in zip(bot[:3],strong[:3]))<=strong[3] for bot in bots)
 low=[min(b[i]-b[3] for b in bots) for i in range(3)];high=[max(b[i]+b[3] for b in bots) for i in range(3)];size=1
 while size<max(b-a+1 for a,b in zip(low,high)):size*=2
 def pointbox(point,lo,size):return sum(max(a-p,0,p-(a+size-1)) for p,a in zip(point,lo))
 def boxpriority(lo,size):return -sum(pointbox(b[:3],lo,size)<=b[3] for b in bots),pointbox((0,0,0),lo,size),size,tuple(lo)
 queue=[boxpriority(low,size)];processed=0;started=time.monotonic()
 while queue:
  negcount,d,size,lo=heapq.heappop(queue);processed+=1
  if size==1:break
  if processed%1000==0 and time.monotonic()-started>70:raise TimeoutError('nanobotoraclebudget')
  half=size//2
  for delta in itertools.product((0,half),repeat=3):child=tuple(a+b for a,b in zip(lo,delta));heapq.heappush(queue,boxpriority(child,half))
 save(23,part1,d,optimal_point=lo,optimal_coverage=-negcount,boxes_processed=processed)
 # Constallations are connected components ofpairManhattandistance<=3.
 points=[tuple(map(int,l.split(','))) for l in raw(25).splitlines()];parent=list(range(len(points)))
 def find(i):
  while parent[i]!=i:parent[i]=parent[parent[i]];i=parent[i]
  return i
 for i,a in enumerate(points):
  for j in range(i):
   if sum(abs(x-y) for x,y in zip(a,points[j]))<=3:parent[find(i)]=find(j)
 save(25,len({find(i) for i in range(len(points))}))
# StageC: canonical and actual combat, independent stable-ID units and two BFS passes.
def combat(text,elf_attack=3,abort_elf=False):
 rows=text.splitlines();floor={(x,y) for y,row in enumerate(rows) for x,c in enumerate(row) if c!='#'};units=[{'pos':(x,y),'side':c,'hp':200} for y,row in enumerate(rows) for x,c in enumerate(row) if c in 'EG'];occupied={u['pos']:i for i,u in enumerate(units)};rounds=0;deaths=0;dirs=[(0,-1),(-1,0),(1,0),(0,1)]
 def adjacent(p):return [(p[0]+dx,p[1]+dy) for dx,dy in dirs]
 def bfs(start,allowed):
  distance={start:0};queue=collections.deque([start])
  while queue:
   p=queue.popleft()
   for q in adjacent(p):
    if q in allowed and q not in distance:distance[q]=distance[p]+1;queue.append(q)
  return distance
 while True:
  assert rounds<10000;order=sorted([i for i,u in enumerate(units) if u['hp']>0],key=lambda i:(units[i]['pos'][1],units[i]['pos'][0]));changed=False
  for i in order:
   u=units[i]
   if u['hp']<=0:continue
   enemies=[j for j,v in enumerate(units) if v['hp']>0 and v['side']!=u['side']]
   if not enemies:return rounds*sum(v['hp'] for v in units if v['hp']>0),deaths,rounds,sum(v['hp'] for v in units if v['hp']>0)
   def targets():return [occupied[p] for p in adjacent(u['pos']) if p in occupied and units[occupied[p]]['side']!=u['side']]
   attack=targets()
   if not attack:
    allowed=floor-set(occupied);allowed.add(u['pos']);inrange={p for j in enemies for p in adjacent(units[j]['pos']) if p in allowed};distance=bfs(u['pos'],allowed);reachable=[p for p in inrange if p in distance]
    if reachable:
     dest=min(reachable,key=lambda p:(distance[p],p[1],p[0]));back=bfs(dest,allowed);step=next(p for p in adjacent(u['pos']) if p in back and back[p]==distance[dest]-1);occupied.pop(u['pos']);u['pos']=step;occupied[step]=i;changed=True
    attack=targets()
   if attack:
    j=min(attack,key=lambda j:(units[j]['hp'],units[j]['pos'][1],units[j]['pos'][0]));v=units[j];v['hp']-=elf_attack if u['side']=='E' else 3;changed=True
    if v['hp']<=0:
     occupied.pop(v['pos'])
     if v['side']=='E':
      deaths+=1
      if abort_elf:return None,deaths,rounds,0
  if not changed:raise ValueError('combat stalemate')
  rounds+=1

def battle(groups,boost=0):
 g=[dict(v,units=v['units'],damage=v['damage']+(boost if v['side']==0 else 0)) for v in groups]
 def ep(i):return g[i]['units']*g[i]['damage']
 def damage(i,j):return 0 if g[i]['kind'] in g[j]['immune'] else ep(i)*(2 if g[i]['kind'] in g[j]['weak'] else 1)
 rounds=0
 while len({v['side'] for v in g if v['units']>0})>1:
  rounds+=1;assert rounds<100000;targeted=set();targets={}
  for i in sorted([i for i,v in enumerate(g) if v['units']>0],key=lambda i:(-ep(i),-g[i]['initiative'])):
   available=[j for j,v in enumerate(g) if v['units']>0 and v['side']!=g[i]['side'] and j not in targeted and damage(i,j)>0]
   if available:j=max(available,key=lambda j:(damage(i,j),ep(j),g[j]['initiative']));targets[i]=j;targeted.add(j)
  killed=0
  for i in sorted(targets,key=lambda i:-g[i]['initiative']):
   j=targets[i]
   if g[i]['units']<=0 or g[j]['units']<=0:continue
   n=min(g[j]['units'],damage(i,j)//g[j]['hp']);g[j]['units']-=n;killed+=n
  if killed==0:return None,sum(v['units'] for v in g),rounds
 winner=next(v['side'] for v in g if v['units']>0);return winner,sum(v['units'] for v in g if v['units']>0),rounds

def parse_armies(immune,infection):
 groups=[]
 for side,text in enumerate((immune,infection)):
  for l in text.splitlines():
   n,hp,extras,attack,kind,initiative=re.fullmatch(r'(\d+) units each with (\d+) hit points(?: \(([^)]+)\))? with an attack that does (\d+) (\w+) damage at initiative (\d+)',l).groups();weak=set();invuln=set()
   for spec in extras.split('; ') if extras else []:
    head,values=spec.split(' to ');(weak if head=='weak' else invuln).update(values.split(', '))
   groups.append(dict(side=side,units=int(n),hp=int(hp),damage=int(attack),kind=kind,initiative=int(initiative),weak=weak,immune=invuln))
 return groups

def stage_c():
 # Anchor simulator on independent published combat outcomes before accepting actual goldens.
 expected={2:27730,3:36334,4:39514,5:27755,6:28944,7:18740};observed={}
 for i,want in expected.items():
  got=combat(COMBAT_EXAMPLES[i])[0];assert got==want,(i,got,want);observed[i]=got
 first=combat(raw(15));started=time.monotonic()
 for power in itertools.count(3):
  outcome=combat(raw(15),power,True)
  if time.monotonic()-started>90:raise TimeoutError('elfboostoraclebudget')
  if outcome[0] is not None:break
 save(15,first[0],outcome[0],canonical_outcomes=observed,minimum_elf_attack=power,part1_rounds=first[2],part1_hp=first[3],part2_rounds=outcome[2],part2_hp=outcome[3])
 groups=parse_armies('17 units each with 5390 hit points (weak to radiation, bludgeoning) with an attack that does 4507 fire damage at initiative 2\n989 units each with 1274 hit points (immune to fire; weak to bludgeoning, slashing) with an attack that does 25 slashing damage at initiative 3','801 units each with 4706 hit points (weak to radiation) with an attack that does 116 bludgeoning damage at initiative 1\n4485 units each with 2961 hit points (immune to radiation; weak to fire, cold) with an attack that does 12 slashing damage at initiative 4');assert battle(groups)[:2]==(1,5216);assert battle(groups,1570)[:2]==(0,51)
 groups=parse_armies(raw(24).split('\n\nInfection:\n', 1)[0].removeprefix('Immune System:\n'),raw(24).split('\n\nInfection:\n', 1)[1]);first=battle(groups);started=time.monotonic()
 for boost in itertools.count(0):
  winner,survivors,rounds=battle(groups,boost)
  if time.monotonic()-started>90:raise TimeoutError('armyboostoraclebudget')
  if winner==0:break
 save(24,first[1],survivors,minimum_boost=boost,part1_winner=first[0],part1_rounds=first[2],part2_rounds=rounds)

COMBAT_EXAMPLES = {2: '#######\n#.G...#\n#...EG#\n#.#.#G#\n#..G#E#\n#.....#\n#######', 3: '#######\n#G..#E#\n#E#E.E#\n#G.##.#\n#...#E#\n#...E.#\n#######', 4: '#######\n#E..EG#\n#.#G.E#\n#E.##E#\n#G..#.#\n#..E#.#\n#######', 5: '#######\n#E.G#.#\n#.#G..#\n#G.#.G#\n#G..#.#\n#...E.#\n#######', 6: '#######\n#.E...#\n#.#..G#\n#.###.#\n#E#G#G#\n#...#G#\n#######', 7: '#########\n#G......#\n#.E.#...#\n#..##..G#\n#...##..#\n#...#...#\n#.G...G.#\n#.....G.#\n#########'}
stage_b()
stage_c()
assert len(OUT) == 25
for day, record in OUT.items():
 assert [record['part1'], record['part2']] == record['source_expected'], day
args.output.write_text(json.dumps(dict(sorted(OUT.items(), key=lambda pair: int(pair[0]))), indent=2) + '\n')
