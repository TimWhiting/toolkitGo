package toolkitGo

import "math/rand"

type MLSystemManager struct{

}
type ArgParser struct{
	Arff,Learner,Evaluation,EvalExtra string
	Verbose,Normalize bool
}

func main(){
	ml := MLSystemManager{};
	args := ArgParser{}
	ml.Run(args);
}

func (ml MLSystemManager)GetLearner(mode1 string, rand rand.Rand)(SupervisedLearner, error){
 return SupervisedLearner{}, nil;
}
func (ml MLSystemManager)Run(args ArgParser)(error){
	return nil;
}