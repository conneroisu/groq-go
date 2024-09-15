library IEEE;

use IEEE.std_logic_1164.all;

entity invg is

  port(
    i_A : in  std_logic;
    o_F : out std_logic
    );

end invg;

architecture dataflow of invg is
begin

  o_F <= not i_A;

end dataflow;
